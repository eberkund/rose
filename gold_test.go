package rose_test

import (
	"flag"
	"testing"

	"github.com/eberkund/rose"
)

var update bool

func init() {
	flag.BoolVar(&update, "update", false, "update golden files")
}

func TestGold_JSONEq(t *testing.T) {
	g := rose.New(t, rose.UpdateFlag(update))
	g.JSONEq("testdata/json_eq.golden.json", `{"foo":123,"bar":"hello world","a":true}`)
}

func TestGold_TOMLEq(t *testing.T) {
	g := rose.New(t, rose.UpdateFlag(update), rose.Prefix("testdata"))
	g.TOMLEq("toml_eq.golden.toml", `
	Age = 25
	Cats = [ "Cauchy", "Plato" ]
	
	Pi = 3.14
	Perfection = [ 6, 28, 496, 8128 ]
	DOB = 1987-07-05T05:45:00Z
	`)
}

func TestGold_Eq(t *testing.T) {
	g := rose.New(t, rose.UpdateFlag(update), rose.Prefix("testdata"))
	g.Eq("text_eq.golden.txt", "Hello\nWorld\n!")
}
