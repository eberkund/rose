package rose

import (
	"flag"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type Gold struct {
	ignoreOrder bool
	flagName    string
	flag        bool
	t           *testing.T
}

type GoldOption func(*Gold)

func IgnoreOrder() GoldOption {
	return func(g *Gold) {
		g.ignoreOrder = true
	}
}

func FlagName(flagName string) GoldOption {
	return func(g *Gold) {
		g.flagName = flagName
	}
}

func New(t *testing.T, options ...GoldOption) *Gold {
	g := &Gold{
		t:           t,
		ignoreOrder: false,
		flagName:    "update",
		flag:        false,
	}
	for _, o := range options {
		o(g)
	}
	flag.BoolVar(&g.flag, g.flagName, false, "update golden files")
	return g
}

func (g *Gold) JSONEq(golden string, reader io.Reader) {
	if g.flag {
		file, err := os.OpenFile(golden, os.O_WRONLY, os.ModeExclusive)
		require.NoError(g.t, err)
		err = formatJSON(reader, file)
		require.NoError(g.t, err)
	}
	got, err := ioutil.ReadAll(reader)
	require.NoError(g.t, err)
	expected, err := ioutil.ReadFile(golden)
	require.NoError(g.t, err)
	require.JSONEq(g.t, string(expected), string(got))
}

func (g *Gold) YAMLEq(golden string, reader io.Reader) {

}
