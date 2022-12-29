package datum

import (
	"io"
)

type DatumStream interface {
	Next() (Datum, error)
}

func StreamToSlice(stream DatumStream) ([]Datum, error) {
	datums := []Datum{}
	for {
		datum, err := stream.Next()
		if err == io.EOF {
			return datums, nil
		} else if err != nil {
			return nil, err
		}
		datums = append(datums, datum)
	}
}
