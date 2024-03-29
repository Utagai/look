package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"

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
	"github.com/utagai/look/config"
	"github.com/utagai/look/config/custom"
	"github.com/utagai/look/data"
	"github.com/utagai/look/datum"
	"github.com/utagai/look/query"
)

// DefaultBufSizeBytes is 100 MB.
const DefaultBufSizeBytes = 100 * (2 << 20)

func main() {
	cfg, err := config.Get()
	if err != nil {
		log.Fatalf("failed to get a configuration: %v", err)
	}
	var src io.Reader = cfg.Source
	if cfg.CustomFields != nil {
		src, err = custom.NewFieldsReader(src, cfg.CustomFields, DefaultBufSizeBytes)
		if err != nil {
			log.Fatalf("failed to create a custom fields reader: %v", err)
		}
	}

	// FIXME: This is dangerous if the source is large.
	bytes, err := io.ReadAll(src)
	if err != nil {
		log.Fatalf("failed to read all the bytes from the source (%q): %v", cfg.Source.Name(), err)
	}

	datums := []datum.Datum{}
	if err := json.Unmarshal(bytes, &datums); err != nil {
		log.Fatalf("failed to unmarshal data into JSON: %v", err)
	}

	var d data.Data
	switch cfg.Backend.Type {
	case config.BackendTypeMemory:
		d = data.NewMemoryData(datums, query.NewLiquidQueryExecutor())
	case config.BackendTypeMongoDB:
		d, err = data.NewMongoDBData(cfg.Backend.MongoDB, "look", cfg.Source.Name(), datums)
		if err != nil {
			log.Fatalf("failed to create the MongoDB backend: %v", err)
		}
	default:
		log.Fatalf("unexpected backend type %q", cfg.Backend.Type)
	}

	initializeGowid(d)
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
			queryStatusTextbox.SetText("Done.", app)
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
