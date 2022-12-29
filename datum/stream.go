package datum

import (
	"io"
)

type Stream interface {
	Next() (Datum, error)
}

func StreamToSlice(stream Stream) ([]Datum, error) {
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
