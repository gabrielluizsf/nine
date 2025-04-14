package json

import "bytes"

// Buffer is an interface that defines a method for converting data into a byte slice.
type Buffer interface {
	Bytes() ([]byte, error)
	Buffer() (*bytes.Buffer, error)
}

func RWBuffer(buf Buffer) (*bytes.Buffer, error) {
	b, err := buf.Bytes()
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(b), nil
}
