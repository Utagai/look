package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/gcla/gowid"
	"github.com/gcla/gowid/examples"
	"github.com/gcla/gowid/widgets/edit"
	"github.com/gcla/gowid/widgets/framed"
	"github.com/gcla/gowid/widgets/holder"
	"github.com/gcla/gowid/widgets/list"
	"github.com/gcla/gowid/widgets/pile"
	"github.com/gcla/gowid/widgets/styled"
	"github.com/gcla/gowid/widgets/text"
	"github.com/gcla/gowid/widgets/vpadding"
	"github.com/utagai/look/data"
	"github.com/utagai/look/datum"
	"github.com/utagai/look/query"
)

type BackendType string

const (
	BackendTypeMemory  = "memory"
	BackendTypeMongoDB = "mongodb"
)

type Config struct {
	Source  *os.File
	Backend struct {
		Type    BackendType
		Memory  bool
		MongoDB string
	}
}

func main() {
	cfg := getConfig()
	// TODO: This is dangerous if the source is large.
	bytes, err := ioutil.ReadAll(cfg.Source)
	if err != nil {
		log.Fatalf("failed to read all the bytes from the source (%q): %v", cfg.Source.Name(), err)
	}

	datums := []datum.Datum{}
	// TODO: We should probably do something more abstracted here than just
	// auto-unmarshalling from JSON.
	if err := json.Unmarshal(bytes, &datums); err != nil {
		log.Fatalf("failed to unmarshal data into JSON: %v", err)
	}

	var d data.Data
	switch cfg.Backend.Type {
	case BackendTypeMemory:
		d = data.NewMemoryData(datums)
	case BackendTypeMongoDB:
		d, err = data.NewMongoDBData(cfg.Backend.MongoDB, "look", cfg.Source.Name(), datums)
		if err != nil {
			log.Fatalf("failed to create the MongoDB backend: %v", err)
		}
	default:
		log.Fatalf("unexpected backend type %q", cfg.Backend.Type)
	}

	initializeGowid(d)
}

func getConfig() *Config {
	sourcePtr := flag.String("source", "", "the source of data")
	mongodbPtr := flag.String("mongodb", "", "specify the MongoDB connection string URI")

	flag.Parse()

	//// Validate.
	if *sourcePtr == "" {
		log.Fatalf("must specify a source of data")
	}

	//// Set onto Config.
	var cfg Config

	// Source
	source := *sourcePtr
	fi := os.Stdin
	var err error
	if source != "-" {
		fi, err = os.Open(source)
		if err != nil {
			log.Fatalf("failed to open source (%q): %v", source, err)
		}
	}

	cfg.Source = fi

	// Backend type.
	cfg.Backend.Type = BackendTypeMemory
	if *mongodbPtr != "" {
		cfg.Backend.Type = BackendTypeMongoDB
		cfg.Backend.MongoDB = *mongodbPtr
	}

	return &cfg
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
		text.StyledContent("look | ", title),
		text.StyledContent("ESC", key),
		text.StringContent(" exits."),
	}

	footerText := styled.New(text.NewFromContent(text.NewContent(footerContent)), foot)

	walker := data.NewDataWalker(d)
	lb := list.NewBounded(walker)
	styledLb := styled.New(lb, body)

	queryTextbox := edit.New(edit.Options{Caption: "Query: "})

	framedQueryTextboxValid := framed.New(queryTextbox, framed.Options{
		Frame: framed.UnicodeFrame,
		Style: gowid.MakeForeground(gowid.ColorGreen),
	})
	framedQueryTextboxInvalid := framed.New(queryTextbox, framed.Options{
		Frame: framed.UnicodeFrame,
		Style: gowid.MakeForeground(gowid.ColorRed),
	})

	queryStatusTextbox := edit.New(edit.Options{Text: "Done."})

	framedQueryStatusTextboxValid := framed.New(queryStatusTextbox, framed.Options{
		Frame: framed.UnicodeFrame,
		Style: gowid.MakeForeground(gowid.ColorGreen),
	})
	framedQueryStatusTextboxInvalid := framed.New(queryStatusTextbox, framed.Options{
		Frame: framed.UnicodeFrame,
		Style: gowid.MakeForeground(gowid.ColorRed),
	})
	queryStatusTextboxHolder := holder.New(framedQueryStatusTextboxValid)

	queryTextboxHolder := holder.New(framedQueryTextboxValid)
	queryTextbox.OnTextSet(gowid.WidgetCallback{
		Name: "on query text change",
		WidgetChangedFunction: func(app gowid.IApp, w gowid.IWidget) {
			newData, err := d.Find(context.Background(), queryTextbox.Text())
			if errors.Is(err, query.ErrUnableToParseQuery) {
				log.Printf("incomplete query: %q", queryTextbox.Text())
				queryTextboxHolder.SetSubWidget(framedQueryTextboxInvalid, app)
				queryStatusTextboxHolder.SetSubWidget(framedQueryStatusTextboxInvalid, app)
				queryStatusTextbox.SetText(err.Error(), app)
				return
			} else if err != nil {
				log.Fatalf("failed to construct the new data: %v", err)
			}
			queryTextboxHolder.SetSubWidget(framedQueryTextboxValid, app)
			queryStatusTextboxHolder.SetSubWidget(framedQueryStatusTextboxValid, app)
			queryStatusTextbox.SetText("Done", app)
			lb.SetWalker(data.NewDataWalker(newData), app)
		},
	})

	view := pile.New([]gowid.IContainerWidget{
		&gowid.ContainerWidget{
			IWidget: vpadding.New(queryTextboxHolder, gowid.VAlignMiddle{}, gowid.RenderFlow{}),
			D:       gowid.RenderFlow{},
		},
		&gowid.ContainerWidget{
			IWidget: styledLb,
			D:       gowid.RenderWithWeight{W: 1},
		},
		&gowid.ContainerWidget{
			IWidget: vpadding.New(queryStatusTextboxHolder, gowid.VAlignMiddle{}, gowid.RenderFlow{}),
			D:       gowid.RenderFlow{},
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
