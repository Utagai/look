package config

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type customFieldsReader struct {
	parseFields     []ParseField
	fieldValueRegex *regexp.Regexp
	jsonPieceChan   chan jsonPiece
	leftoverBytes   []byte
	bufReader       *bufio.Reader
}

func buildRegexFromParseFields(parseFields []ParseField) (*regexp.Regexp, error) {
	var regexpStr strings.Builder
	for _, parseField := range parseFields {
		regexpStr.WriteString(fmt.Sprintf(".*?(%s)", parseField.Regex))
	}

	return regexp.Compile(regexpStr.String())
}

type jsonPiece struct {
	piece string
	err   error
}

// NewCustomFieldsReader creates a new CustomFieldsReader.
func NewCustomFieldsReader(src io.Reader, parseFields []ParseField) (io.Reader, error) {
	if len(parseFields) == 0 {
		// If there are no custom fields, then just treat this as a no-op.
		return src, nil
	}

	bufReader := bufio.NewReader(src)

	regex, err := buildRegexFromParseFields(parseFields)
	if err != nil {
		// TODO: Should we be compiling each of the individual regexes supplied by
		// the user?
		return nil, fmt.Errorf("failed to compile the combined regex: %w", err)
	}
	jsonPieceChan := make(chan jsonPiece, 100) // TODO: Should be a const?
	cfr := &customFieldsReader{
		parseFields:     parseFields,
		fieldValueRegex: regex,
		bufReader:       bufReader,
		jsonPieceChan:   jsonPieceChan,
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
		if !firstLoop && len(json) != 0 {
			r.jsonPieceChan <- jsonPiece{
				piece: ",",
				err:   nil,
			}
		}
		// Even if we got an error, it doesn't matter because JSON could be
		// non-empty. Always be sure to include it.
		r.jsonPieceChan <- jsonPiece{
			piece: json,
			err:   nil,
		}
		if err != nil {
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
			return
		}
		firstLoop = false
	}
}

func (r *customFieldsReader) mapToJSON(matchedTexts []string) string {
	// The parse fields are ordered by their appearance in the regex, so the
	// submatches are 1:1.
	var jsonString strings.Builder
	jsonString.WriteString("{")
	for i := range matchedTexts {
		matchedText := matchedTexts[i]
		parseField := r.parseFields[i]

		jsonString.WriteString(fmt.Sprintf("%q: %s", parseField.FieldName, matchedText))

		if i != len(matchedTexts)-1 {
			jsonString.WriteString(",")
		}
	}

	jsonString.WriteString("}")

	return jsonString.String()
}

func (r *customFieldsReader) getNextJSONDocument() (string, error) {
	var json string
	for {
		line, err := r.bufReader.ReadString('\n')
		if err != nil {
			return "", err
		}

		submatches := r.fieldValueRegex.FindStringSubmatch(line)
		if len(submatches) <= 0 {
			continue
		}
		// The first submatch is going to be the entirety of the regex, but we don't
		// care about that. We care about the groups:
		json = r.mapToJSON(submatches[1:])

		return json, nil
	}
}

func (r *customFieldsReader) Read(p []byte) (int, error) {
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
