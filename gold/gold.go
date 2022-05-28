package gold

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/eberkund/rose/formatting"
	"github.com/pkg/errors"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/spf13/afero"
)

// Testing is a subset of testing.TB that can be reimplemented.
type Testing interface {
	Cleanup(func())
	Error(args ...any)
	Errorf(format string, args ...any)
	Fail()
	FailNow()
	Failed() bool
	Fatal(args ...any)
	Fatalf(format string, args ...any)
	Helper()
	Log(args ...any)
	Logf(format string, args ...any)
	Name() string
	Setenv(key, value string)
	Skip(args ...any)
	SkipNow()
	Skipf(format string, args ...any)
	Skipped() bool
	TempDir() string
}

// Gold makes assertions against golden files.
type Gold struct {
	flag   bool
	fail   func()
	prefix string
	t      Testing
	fs     afero.Fs
}

// New initializes an instance of Gold.
func New(t Testing, options ...Option) *Gold {
	g := &Gold{
		t:      t,
		flag:   false,
		fail:   t.FailNow,
		prefix: "testdata",
		fs:     afero.NewOsFs(),
	}
	for _, o := range options {
		o(g)
	}
	return g
}

func (g *Gold) prependPrefix(path string) string {
	return filepath.Join(g.prefix, path)
}

func (g *Gold) update(filename, actual string, formatter formatting.Formats) error {
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
		return errors.Wrap(err, "could not format or write data to golden file")
	}
	return nil
}

func (g *Gold) assert(goldenPath, given string, formatter formatting.Formats) (string, error) {
	prefixed := g.prependPrefix(goldenPath)
	if err := g.update(prefixed, given, formatter); err != nil {
		return "", err
	}
	var formatted bytes.Buffer
	if err := formatter(strings.NewReader(given), &formatted); err != nil {
		return "", errors.Wrap(err, "could not format given data")
	}
	expected, err := afero.ReadFile(g.fs, prefixed)
	if err != nil {
		return "", errors.Wrap(err, "could not read golden file")
	}
	text, err := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
		FromFile: "Golden",
		A:        difflib.SplitLines(string(expected)),
		ToFile:   "Given",
		B:        difflib.SplitLines(formatted.String()),
		Context:  3,
	})
	if err != nil {
		return "", errors.Wrap(err, "could not produce diff")
	}
	return text, nil
}

func (g *Gold) verify(diff string, err error) {
	g.t.Helper()
	if err != nil {
		g.t.Log(err)
		g.fail()
	} else if diff != "" {
		g.t.Logf("\n%s", diff)
		g.fail()
	}
}
