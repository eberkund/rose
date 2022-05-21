package formatting

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/pelletier/go-toml"
	"github.com/yosssi/gohtml"
	"gopkg.in/yaml.v3"
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

// JSON formats the contents of reader and writes the results to writer.
func JSON(reader io.Reader, writer io.Writer) error {
	decoder := json.NewDecoder(reader)
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "\t")
	return format(encoder, decoder)
}

// HTML formats the contents of reader and writes the results to writer.
func HTML(reader io.Reader, writer io.Writer) error {
	all, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	formatted := gohtml.Format(string(all))
	_, err = io.WriteString(writer, formatted)
	return err
}

// YAML formats the contents of reader and writes the results to writer.
func YAML(reader io.Reader, writer io.Writer) error {
	encoder := yaml.NewEncoder(writer)
	encoder.SetIndent(2)
	decoder := yaml.NewDecoder(reader)
	return format(encoder, decoder)
}

// TOML formats the contents of reader and writes the results to writer.
func TOML(reader io.Reader, writer io.Writer) error {
	encoder := toml.NewEncoder(writer)
	encoder.Order(toml.OrderAlphabetical)
	encoder.Indentation("  ")
	decoder := toml.NewDecoder(reader)
	return format(encoder, decoder)
}

// NoOp copies the reader contents to writer without any changes.
func NoOp(reader io.Reader, writer io.Writer) error {
	_, err := io.Copy(writer, reader)
	return err
}
