package datum

import (
	"io"
	"log"
)

type DatumStream interface {
	Next() (Datum, error)
}

func StreamToSlice(stream DatumStream) ([]Datum, error) {
	datums := []Datum{}
	for {
		datum, err := stream.Next()
		if err == io.EOF {
			log.Println("AYO EOF, RETURNING NUM DATUMS: ", len(datums))
			return datums, nil
		} else if err != nil {
			return nil, err
		}
		log.Println("OK APPENDING")
		datums = append(datums, datum)
	}
}
