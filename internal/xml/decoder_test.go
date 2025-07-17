package xml

import (
	"testing"

	"github.com/i9si-sistemas/assert"
	"github.com/i9si-sistemas/stringx"
)

func TestDecode(t *testing.T) {
	t.Run("SimpleXML", func(t *testing.T) {
		input := `
		<root>
			<name>Gabriel</name>
			<age>30</age>
		</root>`

		r := stringx.NewReader(input)

		got, err := Decode(r)
		assert.NoError(t, err)

		want := map[string]any{
			"name": map[string]any{
				"#text": "Gabriel",
			},
			"age": map[string]any{
				"#text": "30",
			},
		}

		assert.Equal(t, want, got)
	})

	t.Run("InvalidXML", func(t *testing.T) {
		input := `<root><name>Gabriel</age></root>`
		r := stringx.NewReader(input)

		_, err := Decode(r)
		assert.Error(t, err)
	})

	t.Run("EmptyInput", func(t *testing.T) {
		r := stringx.NewReader(stringx.Empty.String())

		_, err := Decode(r)
		assert.Equal(t, err, ErrInvalidXML)
	})
}
