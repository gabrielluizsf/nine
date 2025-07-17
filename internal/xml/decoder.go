package xml

import (
	"encoding/xml"
	"errors"
	"io"

	"github.com/i9si-sistemas/stringx"
)

var ErrInvalidXML = errors.New("invalid xml")

func Decode(r io.Reader) (map[string]any, error) {
	decoder := xml.NewDecoder(r)
	for {
		tok, err := decoder.Token()
		if err != nil {
			break
		}

		if start, ok := tok.(xml.StartElement); ok {
			return parseElement(decoder, start)
		}
	}
	return nil, ErrInvalidXML
}

func parseElement(d *xml.Decoder, start xml.StartElement) (map[string]any, error) {
	type elementFrame struct {
		start  xml.StartElement
		result map[string]any
	}

	stack := []elementFrame{{start: start, result: map[string]any{}}}

	for {
		token, err := d.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		switch tok := token.(type) {
		case xml.StartElement:
			stack = append(stack, elementFrame{start: tok, result: map[string]any{}})

		case xml.CharData:
			text := stringx.String(string(tok)).Trim(stringx.Space.String()).
			Trim(stringx.Tab.String()).Trim(stringx.NewLine.String()).String()
			if text != "" {
				top := &stack[len(stack)-1]
				top.result["#text"] = text
			}

		case xml.EndElement:
			if len(stack) < 2 {
				return stack[0].result, nil
			}

			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			parent := &stack[len(stack)-1]

			name := top.start.Name.Local
			if existing, ok := parent.result[name]; ok {
				switch v := existing.(type) {
				case []any:
					parent.result[name] = append(v, top.result)
				default:
					parent.result[name] = []any{v, top.result}
				}
			} else {
				parent.result[name] = top.result
			}
		}
	}

	return stack[0].result, nil
}
