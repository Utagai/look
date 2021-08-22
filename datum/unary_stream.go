package datum

import (
	"io"
)

// UnaryStream is a datum stream that returns only a single datum, and then
// reports EOF.
type UnaryStream struct {
	datum         Datum
	datumReturned bool
}

// NewUnaryStream is a constructor for UnaryStream.
func NewUnaryStream(datum Datum) *UnaryStream {
	return &UnaryStream{
		datum:         datum,
		datumReturned: false,
	}
}

// Next implements DatumStream.
func (d *UnaryStream) Next() (Datum, error) {
	if d.datumReturned {
		return nil, io.EOF
	}

	d.datumReturned = true
	return d.datum, nil
}
