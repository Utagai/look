package custom

import (
	"bufio"
	"bytes"
	"errors"
	"io"
)

const pieceChanSize = 100

var (
	errJSONFinished = errors.New("source JSON input finished")
)

type fieldsReader struct {
	fields                *Fields
	jsonPieceChan         chan jsonPiece
	buf                   *bytes.Buffer
	approxMaxBufSizeBytes uint
	src                   *bufio.Reader
}

type jsonPiece struct {
	piece []byte
	err   error
}

// NewFieldsReader creates a new fieldsReader from a source reader and the
// custom fields. Note that the max buf size argument is an approximation -- the
// implementation may exceed it at most by the size of a single JSON document
// created from an input line.
func NewFieldsReader(src io.Reader, fields *Fields, approxMaxBufSizeBytes uint) (io.Reader, error) {
	if fields == nil {
		// If there are no custom fields, then just treat this as a no-op.
		return src, nil
	}

	bufReader := bufio.NewReader(src)

	// NOTE: Using a channel may not be the fastest option. For now though, it
	// keeps the code simpler and cleaner.
	jsonPieceChan := make(chan jsonPiece, pieceChanSize)
	cfr := &fieldsReader{
		fields:                fields,
		src:                   bufReader,
		jsonPieceChan:         jsonPieceChan,
		buf:                   bytes.NewBuffer(make([]byte, 0, approxMaxBufSizeBytes)),
		approxMaxBufSizeBytes: approxMaxBufSizeBytes,
	}

	go cfr.generateJSON()

	return cfr, nil
}

func (r *fieldsReader) generateJSON() {
	// All JSON streams must start with the opening array bracket.
	r.jsonPieceChan <- jsonPiece{
		piece: []byte("["),
		err:   nil,
	}

	firstLoop := true
	for {
		json, err := r.getNextJSONDocument()
		if err != nil {
			r.finalizeJSON(err)
			return
		}

		if !firstLoop {
			r.jsonPieceChan <- jsonPiece{
				piece: []byte(","),
				err:   nil,
			}
		}
		firstLoop = false

		r.jsonPieceChan <- jsonPiece{
			piece: json,
			err:   nil,
		}
	}
}

func (r *fieldsReader) getNextJSONDocument() ([]byte, error) {
	for {
		line, err := r.src.ReadString('\n')
		if err != nil {
			return nil, err
		}

		json, err := r.fields.ToJSON(line)
		if err == ErrNoMatch {
			continue // If this line doesn't match, don't error, just continue to the next line.
		} else if err != nil {
			return nil, err
		}

		return json, nil
	}
}

func (r *fieldsReader) finalizeJSON(err error) {
	// End of JSON, regardless of EOF or more substantial error, so we have a
	// few things to do:
	//   * Finish the JSON by closing the array with ']'.
	//   * Return the actual error so the consumer can determine what it wants
	//   to do with it.
	//   * Close the channel, so we can signal that we are done.
	//   * Return and clean up the goroutine.
	r.jsonPieceChan <- jsonPiece{
		piece: []byte("]"),
		err:   nil,
	}

	r.jsonPieceChan <- jsonPiece{
		piece: []byte{},
		err:   err,
	}

	close(r.jsonPieceChan)
}

func (r *fieldsReader) hydrateBuffer() error {
	// Only hydrate if the buffer is empty.
	if r.buf.Len() > 0 {
		return nil
	}

	// Note that technically, this can exceed maxBufSizeBytes, but only by a
	// single JSON piece's size.
	for r.buf.Len() < int(r.approxMaxBufSizeBytes) {
		jsonRes, ok := <-r.jsonPieceChan
		if !ok {
			// If the channel is closed & exhausted, then there's nothing to hydrate
			// with. Just return errJSONFinished:
			return errJSONFinished
		}

		if err := jsonRes.err; err != nil {
			if err == io.EOF {
				return errJSONFinished
			}
			return err
		}

		_, err := r.buf.Write(jsonRes.piece)
		if err != nil {
			// bytes.Buffer only ever returns err in the case of ErrTooLarge, in which
			// case we've messed up because we should stay below maxBufSize.
			// Technically, this can happen if someone passes in a GIGANTIC line +
			// GIGANTIC regex for parsing out a GIGANTIC JSON document such that a
			// single piece could actually be at the scale of multiple gigabytes...
			// But I don't intend to support that use case.
			panic(err)
		}
	}

	return nil
}

func (r *fieldsReader) Read(p []byte) (int, error) {
	if err := r.hydrateBuffer(); err != nil && err != errJSONFinished {
		return 0, err
	}

	return r.buf.Read(p)
}
