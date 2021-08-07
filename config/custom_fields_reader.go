package config

import (
	"bufio"
	"io"
)

type customFieldsReader struct {
	customFields  *CustomFields
	jsonPieceChan chan jsonPiece
	leftoverBytes []byte
	bufReader     *bufio.Reader
}

type jsonPiece struct {
	piece string
	err   error
}

// NewCustomFieldsReader creates a new CustomFieldsReader.
// TODO: We should return *customFieldsReader.
func NewCustomFieldsReader(src io.Reader, customFields *CustomFields) (io.Reader, error) {
	if customFields == nil {
		// If there are no custom fields, then just treat this as a no-op.
		return src, nil
	}

	bufReader := bufio.NewReader(src)

	// NOTE: Using a channel may not be the fastest option. For now though, it
	// keeps the code simpler and cleaner.
	jsonPieceChan := make(chan jsonPiece, 100) // TODO: Should be a const?
	cfr := &customFieldsReader{
		customFields:  customFields,
		bufReader:     bufReader,
		jsonPieceChan: jsonPieceChan,
	}

	go cfr.generateJSON()

	return cfr, nil
}

func (r *customFieldsReader) generateJSON() {
	// All JSON streams must start with the opening array bracket.
	r.jsonPieceChan <- jsonPiece{
		piece: "[",
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
				piece: ",",
				err:   nil,
			}
		}
		firstLoop = false

		r.jsonPieceChan <- jsonPiece{
			piece: string(json),
			err:   nil,
		}
	}
}

func (r *customFieldsReader) getNextJSONDocument() ([]byte, error) {
	for {
		line, err := r.bufReader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		json, err := r.customFields.ToJSON(line)
		if err == ErrNoMatch {
			continue // If this line doesn't match, don't error, just continue to the next line.
		} else if err != nil {
			return nil, err
		}

		return json, nil
	}
}

func (r *customFieldsReader) finalizeJSON(err error) {
	// End of JSON, regardless of EOF or more substantial error, so we have a
	// few things to do:
	//   * Finish the JSON by closing the array with ']'.
	//   * Return the actual error so the consumer can determine what it wants
	//   to do with it.
	//   * Close the channel, so we can signal that we are done.
	//   * Return and clean up the goroutine.
	r.jsonPieceChan <- jsonPiece{
		piece: "]",
		err:   nil,
	}

	r.jsonPieceChan <- jsonPiece{
		piece: "",
		err:   err,
	}

	close(r.jsonPieceChan)
}

func (r *customFieldsReader) Read(p []byte) (int, error) {
	// TODO: Can we write this simpler by using bytes.Buffer?
	totalNumBytesCopied := 0
	if len(r.leftoverBytes) > 0 {
		numBytesCopied := copy(p, r.leftoverBytes)
		totalNumBytesCopied += numBytesCopied
		r.leftoverBytes = r.leftoverBytes[numBytesCopied:]
		if len(r.leftoverBytes) >= len(p) {
			return numBytesCopied, nil
		}
	}

	for {
		jsonRes, ok := <-r.jsonPieceChan
		if !ok {
			// If the channel is closed & exhausted, always return EOF:
			return totalNumBytesCopied, io.EOF
		}
		if err := jsonRes.err; err != nil {
			return totalNumBytesCopied, err
		}

		jsonBytes := []byte(jsonRes.piece)
		numBytesCopied := copy(p[totalNumBytesCopied:], jsonBytes)
		totalNumBytesCopied += numBytesCopied
		r.leftoverBytes = jsonBytes[numBytesCopied:]
		if totalNumBytesCopied >= len(p) {
			return totalNumBytesCopied, nil
		}
	}
}
