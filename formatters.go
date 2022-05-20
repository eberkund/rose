package rose

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/pelletier/go-toml"
	"github.com/yosssi/gohtml"

	"gopkg.in/yaml.v3"
)

type Formats func(reader io.Reader, writer io.Writer) error

type Encoder interface {
	Encode(v interface{}) error
}

type Decoder interface {
	Decode(v interface{}) error
}

func format(encoder Encoder, decoder Decoder) error {
	var decoded interface{}
	err := decoder.Decode(&decoded)
	if err != nil {
		return err
	}
	return encoder.Encode(decoded)
}

func formatJSON(reader io.Reader, writer io.Writer) error {
	decoder := json.NewDecoder(reader)
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "\t")
	return format(encoder, decoder)
}

func formatHTML(reader io.Reader, writer io.Writer) error {
	all, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	formatted := gohtml.Format(string(all))
	_, err = io.WriteString(writer, formatted)
	return err
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

func formatNoop(reader io.Reader, writer io.Writer) error {
	_, err := io.Copy(writer, reader)
	return err
}
