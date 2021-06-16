package main

import (
	"github.com/gcla/gowid"
	"github.com/gcla/gowid/examples"
	"github.com/gcla/gowid/widgets/list"
	"github.com/gcla/gowid/widgets/pile"
	"github.com/gcla/gowid/widgets/styled"
	"github.com/gcla/gowid/widgets/text"
	"github.com/utagai/look/data"
)

func main() {
	initializeGowid()
}

func initializeGowid() {
	palette := gowid.Palette{
		"title": gowid.MakePaletteEntry(gowid.ColorWhite, gowid.ColorBlack),
		"key":   gowid.MakePaletteEntry(gowid.ColorCyan, gowid.ColorBlack),
		"foot":  gowid.MakePaletteEntry(gowid.ColorWhite, gowid.ColorBlack),
		"body":  gowid.MakePaletteEntry(gowid.ColorBlack, gowid.ColorCyan),
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

	datums := make([]data.Datum, 100)
	for i := range datums {
		datums[i] = map[string]interface{}{
			"foo": "bar",
			"a":   i,
		}
	}
	memData := data.NewMemoryData(datums)
	walker := data.NewDataWalker(memData)
	lb := list.NewBounded(walker)
	styledLb := styled.New(lb, body)

	view := pile.New([]gowid.IContainerWidget{
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
