package gold

import (
	"path"

	"github.com/spf13/afero"
)

// Option is the signature for functions that configure Gold.
type Option func(g *Gold)

// WithFailAfter configures whether to stop test execution upon failure.
func WithFailAfter() Option {
	return func(g *Gold) {
		g.fail = g.t.FailNow
	}
}

// WithFS sets the internal filesystem used for assertions.
func WithFS(fs afero.Fs) Option {
	return func(g *Gold) {
		g.fs = fs
	}
}

// WithFlag sets the formatting option.
func WithFlag(flag bool) Option {
	return func(g *Gold) {
		g.flag = flag
	}
}

// WithPrefix sets the folder prefix for golden files.
func WithPrefix(elems ...string) Option {
	return func(g *Gold) {
		g.prefix = path.Join(elems...)
	}
}
