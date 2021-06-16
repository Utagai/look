package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/gcla/gowid"
	"github.com/gcla/gowid/examples"
	"github.com/gcla/gowid/widgets/divider"
	"github.com/gcla/gowid/widgets/edit"
	"github.com/gcla/gowid/widgets/list"
	"github.com/gcla/gowid/widgets/pile"
	"github.com/gcla/gowid/widgets/styled"
	"github.com/gcla/gowid/widgets/text"
	"github.com/utagai/look/data"
)

type Config struct {
	source *os.File
}

func main() {
	cfg := getConfig()
	bytes, err := ioutil.ReadAll(cfg.source)
	if err != nil {
		log.Fatalf("failed to read all the bytes from the source (%q): %v", cfg.source.Name(), err)
	}

	datums := []data.Datum{}
	if err := json.Unmarshal(bytes, &datums); err != nil {
		log.Fatalf("failed to unmarshal data into JSON: %v", err)
	}

	initializeGowid(data.NewMemoryData(datums))
}

func getConfig() *Config {
	sourcePtr := flag.String("source", "", "the source of data")

	flag.Parse()

	if *sourcePtr == "" {
		log.Fatalf("must specify a source of data")
	}

	source := *sourcePtr
	fi := os.Stdin
	var err error
	if source != "-" {
		fi, err = os.Open(source)
		if err != nil {
			log.Fatalf("failed to open source (%q): %v", source, err)
		}
	}

	return &Config{
		source: fi,
	}
}

func initializeGowid(d data.Data) {
	palette := gowid.Palette{
		"title": gowid.MakePaletteEntry(gowid.ColorWhite, gowid.ColorBlack),
		"key":   gowid.MakePaletteEntry(gowid.ColorCyan, gowid.ColorBlack),
		"foot":  gowid.MakePaletteEntry(gowid.ColorWhite, gowid.ColorBlack),
		"body":  gowid.MakePaletteEntry(gowid.ColorWhite, gowid.ColorNone),
		"fbody": gowid.MakePaletteEntry(gowid.ColorWhite, gowid.ColorBlack),
	}

	key := gowid.MakePaletteRef("key")
	foot := gowid.MakePaletteRef("foot")
	title := gowid.MakePaletteRef("title")
	body := gowid.MakePaletteRef("body")

	footerContent := []text.ContentSegment{
		text.StyledContent("Data Viewer", title),
		text.StringContent("    "),
		text.StyledContent("UP", key),
		text.StringContent(", "),
		text.StyledContent("DOWN", key),
		text.StringContent(", "),
		text.StyledContent("PAGE_UP", key),
		text.StringContent(", "),
		text.StyledContent("PAGE_DOWN", key),
		text.StringContent(", "),
		text.StyledContent("HOME", key),
		text.StringContent(", "),
		text.StyledContent("CTRL-L", key),
		text.StringContent(" move view  "),
		text.StyledContent("Q", key),
		text.StringContent(" exits. Try the mouse wheel."),
	}

	footerText := styled.New(text.NewFromContent(text.NewContent(footerContent)), foot)

	walker := data.NewDataWalker(d)
	lb := list.NewBounded(walker)
	styledLb := styled.New(lb, body)

	textbox := edit.New(edit.Options{Caption: "Query: "})
	textbox.OnTextSet(gowid.WidgetCallback{
		Name: "on query text change",
		WidgetChangedFunction: func(app gowid.IApp, w gowid.IWidget) {
			lb.SetWalker(data.NewDataWalker(d.Find(textbox.Text())), app)
		},
	})

	view := pile.New([]gowid.IContainerWidget{
		&gowid.ContainerWidget{
			IWidget: textbox,
			D:       gowid.RenderFlow{},
		},
		&gowid.ContainerWidget{
			IWidget: divider.NewAscii(),
			D:       gowid.RenderFlow{},
		},
		&gowid.ContainerWidget{
			IWidget: styledLb,
			D:       gowid.RenderWithWeight{W: 1},
		},
		&gowid.ContainerWidget{
			IWidget: footerText,
			D:       gowid.RenderFlow{},
		},
	})

	app, err := gowid.NewApp(gowid.AppArgs{
		View:    view,
		Palette: &palette,
	})
	examples.ExitOnErr(err)

	app.SimpleMainLoop()
}
