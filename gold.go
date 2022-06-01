package rose

import (
	"bytes"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

// Gold makes assertions against golden files.
type Gold struct {
	flag   bool
	prefix string
	t      *testing.T
	fs     afero.Fs
}

// New initializes a Gold.
func New(t *testing.T, options ...GoldOption) *Gold {
	g := &Gold{
		flag:   false,
		t:      t,
		prefix: "testdata",
		fs:     afero.NewOsFs(),
	}
	for _, o := range options {
		o(g)
	}
	return g
}

// GoldOption is a method to configure initialization options.
type GoldOption func(g *Gold)

// WithFS sets the internal filesystem used for assertions.
func WithFS(fs afero.Fs) GoldOption {
	return func(g *Gold) {
		g.fs = fs
	}
}

// WithFlag sets the formatting option for a new instance of Gold.
func WithFlag(flag bool) GoldOption {
	return func(g *Gold) {
		g.flag = flag
	}
}

// WithPrefix sets the folder prefix for golden files.
func WithPrefix(elems ...string) GoldOption {
	return func(g *Gold) {
		g.prefix = path.Join(elems...)
	}
}

func (g *Gold) prependPrefix(path string) string {
	return filepath.Join(g.prefix, path)
}

func (g *Gold) genericEQ(goldenPath, actual string, formatter Formats) {
	withPrefix := g.prependPrefix(goldenPath)
	if g.flag {
		err := g.fs.MkdirAll(filepath.Dir(withPrefix), 0o750)
		require.NoError(g.t, err, "could not create directory which holds golden file")

		create, err := g.fs.Create(withPrefix)
		require.NoError(g.t, err, "could not create golden file")

		err = create.Close()
		require.NoError(g.t, err, "could not close golden file after creating it")

		file, err := g.fs.OpenFile(withPrefix, os.O_WRONLY, os.ModeExclusive)
		require.NoError(g.t, err, "error opening golden file for writing")

		err = formatter(strings.NewReader(actual), file)
		require.NoError(g.t, err, "error formatting or writing input data to golden file")
	}

	var formatted bytes.Buffer
	err := formatter(strings.NewReader(actual), &formatted)
	require.NoError(g.t, err, "error formatting input data")

	expected, err := afero.ReadFile(g.fs, withPrefix)
	require.NoError(g.t, err, "error reading golden file")
	require.Equal(g.t, string(expected), formatted.String(), "input data did not match golden file")
}
