package datum

import (
	"io"
)

type DatumSliceStream struct {
	datums []Datum
	i      int
}

func NewDatumSliceStream(datums []Datum) *DatumSliceStream {
	return &DatumSliceStream{
		datums: datums,
		i:      0,
	}
}

func (d *DatumSliceStream) Next() (Datum, error) {
	if d.i >= len(d.datums) {
		return nil, io.EOF
	}

	datum := d.datums[d.i]
	d.i++
	return datum, nil
}
