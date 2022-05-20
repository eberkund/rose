package rose

import (
	"bytes"
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
	prefix string
	t      *testing.T
}

// New initializes a Gold.
func New(t *testing.T, options ...GoldOption) *Gold {
	g := &Gold{
		flag: false,
		t:    t,
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

// Prefix sets the folder prefix for golden files.
func Prefix(prefix string) GoldOption {
	return func(g *Gold) {
		g.prefix = prefix
	}
}

func (g *Gold) addPrefix(path string) string {
	return filepath.Join(g.prefix, path)
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
