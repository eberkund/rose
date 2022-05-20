package rose

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Gold makes assertions against golden files.
type Gold struct {
	flag   bool
	format bool
	prefix string
	t      *testing.T
}

// New initializes a Gold.
func New(t *testing.T, options ...GoldOption) *Gold {
	g := &Gold{
		flag:   false,
		format: true,
		t:      t,
	}
	for _, o := range options {
		o(g)
	}
	return g
}

// GoldOption is a method to configure initialization options.
type GoldOption func(*Gold)

// UpdateFlag sets the formatting option for a new instance of Gold.
func UpdateFlag(flag bool) GoldOption {
	return func(g *Gold) {
		g.flag = flag
	}
}

// Format sets the formatting option for a new instance of Gold.
func Format(format bool) GoldOption {
	return func(g *Gold) {
		g.format = format
	}
}

// Prefix sets the folder prefix for golden files.
func Prefix(prefix string) GoldOption {
	return func(g *Gold) {
		g.prefix = prefix
	}
}

func (g *Gold) addPrefix(path string) string {
	return filepath.Join(g.prefix, path)
}

// JSONEq compares XML to golden file.
func (g *Gold) JSONEq(goldenPath, actual string) {
	g.genericEQ(g.addPrefix(goldenPath), actual, formatJSON)
}

// XMLEq compares XML to golden file.
func (g *Gold) XMLEq(goldenPath, actual string) {
	g.genericEQ(g.addPrefix(goldenPath), actual, formatXML)
}

// TOMLEq compares TOML to golden file.
func (g *Gold) TOMLEq(goldenPath, actual string) {
	g.genericEQ(g.addPrefix(goldenPath), actual, formatTOML)
}

// YAMLEq compares YAML to golden file.
func (g *Gold) YAMLEq(goldenPath, actual string) {
	g.genericEQ(g.addPrefix(goldenPath), actual, formatYAML)
}

// Eq compares string to golden file.
func (g *Gold) Eq(goldenPath, actual string) {
	noopFormatter := func(reader io.Reader, writer io.Writer) error {
		_, err := io.Copy(writer, reader)
		return err
	}
	g.genericEQ(g.addPrefix(goldenPath), actual, noopFormatter)
}

func (g *Gold) genericEQ(goldenPath, actual string, formatter Formats) {
	// Update file if flag is set
	if g.flag {
		file, err := os.OpenFile(goldenPath, os.O_WRONLY, os.ModeExclusive)
		require.NoError(g.t, err)
		err = formatter(strings.NewReader(actual), file)
		require.NoError(g.t, err)
	}

	// Normalize input before comparing
	var formatted bytes.Buffer
	err := formatter(strings.NewReader(actual), &formatted)
	require.NoError(g.t, err)

	// Compare golden to actual
	expected, err := ioutil.ReadFile(goldenPath)
	require.NoError(g.t, err)
	require.Equal(g.t, string(expected), formatted.String())
}
