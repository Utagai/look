package data

import (
	"context"
	"errors"
	"log"

	"github.com/gcla/gowid"
	"github.com/gcla/gowid/widgets/hpadding"
	"github.com/gcla/gowid/widgets/list"
	"github.com/gcla/gowid/widgets/palettemap"
	"github.com/gcla/gowid/widgets/selectable"
	"github.com/gcla/gowid/widgets/text"
)

// DataWalker is a gowid List walker for Data.
// This implementation avoids potential memory issues that
// list.SimpleListWalker would suffer from if the amount data of was very
// large. This implementation gives control of the implementation and its
// memory characteristics to the underlying Data implementation, so it could
// still suffer from memory issue if the Data implementation is not implemented
// well.
type DataWalker struct {
	data  Data
	focus list.IWalkerPosition
	ctx   context.Context
}

var _ list.IBoundedWalker = (*DataWalker)(nil)
var _ list.IWalkerHome = &DataWalker{}
var _ list.IWalkerEnd = &DataWalker{}

func NewDataWalker(data Data) *DataWalker {
	return &DataWalker{
		data:  data,
		ctx:   context.Background(),
		focus: list.ListPos(0),
	}
}

// First implements the list.IWalkerHome interface.
// TODO: Our variable name should probably be something other than 'f'.
func (f *DataWalker) First() list.IWalkerPosition {
	return list.ListPos(0)
}

// Last implements the list.IWalkerEnd interface.
func (f *DataWalker) Last() list.IWalkerPosition {
	length, err := f.data.Length(f.ctx)
	if err != nil {
		log.Fatalf("failed to fetch the length for data: %v", err)
	}
	return list.ListPos(length - 1)
}

func createWidgetFor(datum Datum) gowid.IWidget {
	var res gowid.IWidget
	res = selectable.New(
		palettemap.New(
			hpadding.New(
				text.NewFromContent(
					text.NewContent([]text.ContentSegment{
						text.StyledContent(datum.String(), gowid.MakePaletteRef("body")),
					}),
				),
				gowid.HAlignRight{}, gowid.RenderFlow{},
			),
			palettemap.Map{"body": "fbody"},
			palettemap.Map{},
		),
	)
	return res
}

// At implements the list.IBoundedWalker interface.
func (f *DataWalker) At(pos list.IWalkerPosition) gowid.IWidget {
	index := int(pos.(list.ListPos))
	datum, err := f.data.At(f.ctx, index)
	if errors.Is(err, ErrOutOfBounds) {
		return nil
	} else if err != nil {
		log.Printf("Failed to get datum at index %d: %v", index, err)
		return nil
	}

	return createWidgetFor(datum)
}

// Focus implements the list.IBoundedWalker interface.
func (f *DataWalker) Focus() list.IWalkerPosition {
	return f.focus
}

// SetFocus implements the list.IBoundedWalker interface.
func (f *DataWalker) SetFocus(pos list.IWalkerPosition, app gowid.IApp) {
	f.focus = pos
}

// Next implements the list.IBoundedWalker interface.
func (f *DataWalker) Next(ipos list.IWalkerPosition) list.IWalkerPosition {
	pos := ipos.(list.ListPos)
	length, err := f.data.Length(f.ctx)
	if err != nil {
		log.Fatalf("failed to get the length for data: %v", err)
	}
	if int(pos) == length {
		return list.ListPos(-1)
	} else {
		return pos + 1
	}
}

// Previous implements the list.IBoundedWalker interface.
func (f *DataWalker) Previous(ipos list.IWalkerPosition) list.IWalkerPosition {
	pos := ipos.(list.ListPos)
	if int(pos) == 0 {
		return list.ListPos(-1)
	} else {
		return pos - 1
	}
}

// Length implements the list.IBoundedWalker interface.
func (f *DataWalker) Length() int {
	length, err := f.data.Length(f.ctx)
	if err != nil {
		log.Fatalf("failed to get the length for data: %v", err)
	}
	return length
}
