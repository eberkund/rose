package rose

import (
	"flag"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var update bool

func init() {
	flag.BoolVar(&update, "update", false, "update golden files")
}

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
	}
	for _, o := range options {
		o(g)
	}
	return g
}

func (g *Gold) JSONEq(golden, actual string) {
	if update {
		file, err := os.OpenFile(golden, os.O_WRONLY, os.ModeExclusive)
		require.NoError(g.t, err)
		err = formatJSON(strings.NewReader(actual), file)
		require.NoError(g.t, err)
	}
	expected, err := ioutil.ReadFile(golden)
	require.NoError(g.t, err)
	require.JSONEq(g.t, string(expected), actual)
}

func (g *Gold) YAMLEq(golden string, reader io.Reader) {

}
