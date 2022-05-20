package rose

import (
	"encoding/json"
	"encoding/xml"
	"io"

	"github.com/pelletier/go-toml"

	"gopkg.in/yaml.v3"
)

type Formats func(reader io.Reader, writer io.Writer) error

type Encoder interface {
	Encode(v any) error
}

type Decoder interface {
	Decode(v any) error
}

func format(encoder Encoder, decoder Decoder) error {
	var decoded map[string]interface{}
	err := decoder.Decode(&decoded)
	if err != nil {
		return err
	}
	return encoder.Encode(decoded)
}

func formatJSON(reader io.Reader, writer io.Writer) error {
	decoder := json.NewDecoder(reader)
	encoder := json.NewEncoder(writer)
	return format(encoder, decoder)
}

func formatXML(reader io.Reader, writer io.Writer) error {
	encoder := xml.NewEncoder(writer)
	encoder.Indent("", "\t")
	decoder := xml.NewDecoder(reader)
	return format(encoder, decoder)
}

func formatYAML(reader io.Reader, writer io.Writer) error {
	encoder := yaml.NewEncoder(writer)
	encoder.SetIndent(2)
	decoder := yaml.NewDecoder(reader)
	return format(encoder, decoder)
}

func formatTOML(reader io.Reader, writer io.Writer) error {
	encoder := toml.NewEncoder(writer)
	encoder.Order(toml.OrderAlphabetical)
	encoder.Indentation("  ")
	decoder := toml.NewDecoder(reader)
	return format(encoder, decoder)
}
