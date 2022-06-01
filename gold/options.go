package gold

import (
	"path"

	"github.com/spf13/afero"
)

// Option is the signature for functions that configure Gold.
type Option func(g *Gold)

// WithFatal configures whether file differences will be fail or errors.
func WithFatal(fatal bool) Option {
	return func(g *Gold) {
		if fatal {
			g.fail = g.t.FailNow
		} else {
			g.fail = g.t.Fail
		}
	}
}

// WithFS sets the internal filesystem used for assertions.
func WithFS(fs afero.Fs) Option {
	return func(g *Gold) {
		g.fs = fs
	}
}

// WithFlag sets the formatting option for a new instance of Gold.
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
