package rose

import (
	"bytes"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/spf13/afero"
)

// Gold makes assertions against golden files.
type Gold struct {
	flag   bool
	fatal  bool
	prefix string
	t      *testing.T
	fs     afero.Fs
}

// New initializes a Gold.
func New(t *testing.T, options ...GoldOption) *Gold {
	g := &Gold{
		flag:   false,
		fatal:  true,
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

// WithFatal configures whether file differences will be fatal or errors.
func WithFatal(fatal bool) GoldOption {
	return func(g *Gold) {
		g.fatal = fatal
	}
}

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

func (g *Gold) update(filename, actual string, formatter Formats) error {
	if !g.flag {
		return nil
	}
	err := g.fs.MkdirAll(filepath.Dir(filename), 0o750)
	if err != nil {
		return errors.Wrap(err, "could not create directory which holds golden file")
	}
	create, err := g.fs.Create(filename)
	if err != nil {
		return errors.Wrap(err, "could not create golden file")
	}
	err = create.Close()
	if err != nil {
		return errors.Wrap(err, "could not close golden file after creating it")
	}
	file, err := g.fs.OpenFile(filename, os.O_WRONLY, os.ModeExclusive)
	if err != nil {
		return errors.Wrap(err, "could not open golden file for writing")
	}
	err = formatter(strings.NewReader(actual), file)
	if err != nil {
		return errors.Wrap(err, "could not format or write inputData data to golden file")
	}
	return nil
}

func (g *Gold) assert(goldenPath, actual string, formatter Formats, msgAndArgs ...interface{}) (string, error) {
	prefixed := g.prependPrefix(goldenPath)
	if err := g.update(prefixed, actual, formatter); err != nil {
		return "", err
	}
	var formatted bytes.Buffer
	if err := formatter(strings.NewReader(actual), &formatted); err != nil {
		return "", errors.Wrap(err, "could not format inputData data")
	}
	expected, err := afero.ReadFile(g.fs, prefixed)
	if err != nil {
		return "", errors.Wrap(err, "could not read golden file")
	}
	text, err := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
		FromFile: "Golden",
		A:        difflib.SplitLines(string(expected)),
		ToFile:   "Provided",
		B:        difflib.SplitLines(formatted.String()),
		Context:  3,
	})
	if err != nil {
		return "", errors.Wrap(err, "could not produce diff")
	}
	return text, nil
}

func (g *Gold) handleError(diff string, err error) {
	if err != nil {
		if g.fatal {
			g.t.Fatal(err)
		} else {
			g.t.Error(err)
		}
	}
	if diff != "" {
		if g.fatal {
			g.t.Fatalf("\n%s", diff)
		} else {
			g.t.Errorf("\n%s", diff)
		}
	}
}
