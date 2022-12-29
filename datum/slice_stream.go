package datum

import (
	"io"
)

type SliceStream struct {
	datums []Datum
	i      int
}

func NewSliceStream(datums []Datum) *SliceStream {
	return &SliceStream{
		datums: datums,
		i:      0,
	}
}

func (d *SliceStream) Next() (Datum, error) {
	if d.i >= len(d.datums) {
		return nil, io.EOF
	}

	datum := d.datums[d.i]
	d.i++
	return datum, nil
}
