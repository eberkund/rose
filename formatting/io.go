package formatting

import (
	"io"
)

// Formats is an alias for function signature that reads from reader, formats it and writes to writer.
type Formats func(reader io.Reader, writer io.Writer) error

type encoder interface {
	Encode(v interface{}) error
}

type decoder interface {
	Decode(v interface{}) error
}

func format(encoder encoder, decoder decoder) error {
	var decoded interface{}
	err := decoder.Decode(&decoded)
	if err != nil {
		return err
	}
	return encoder.Encode(decoded)
}
