package rose_test

import (
	"flag"
	"testing"

	"github.com/eberkund/rose"
)

var update = flag.Bool("update", false, "update golden files")

func TestExistingFiles(t *testing.T) {
	testcases := map[string]struct {
		input      string
		goldenFile string
		test       func(g *rose.Gold) func(string, string)
	}{
		"json": {
			input:      `{"foo":123,"bar":"hello world","a":true}`,
			goldenFile: "json_eq.golden.json",
			test: func(gold *rose.Gold) func(string, string) {
				return gold.JSONEq
			},
		},
		"text": {
			goldenFile: "text_eq.golden.txt",
			input:      "Hello\nWorld\n!",
			test: func(gold *rose.Gold) func(string, string) {
				return gold.Eq
			},
		},
		"html": {
			goldenFile: "xml_eq.golden.toml",
			input:      `<fruits><apple/><banana/></fruits>`,
			test: func(gold *rose.Gold) func(string, string) {
				return gold.HTMLEq
			},
		},
		"toml": {
			goldenFile: "toml_eq.golden.toml",
			input: `
Age = 25
Cats = [ "Cauchy", "Plato" ]

Pi = 3.14
Perfection = [ 6, 28, 496, 8128 ]
DOB = 1987-07-05T05:45:00Z
`,
			test: func(gold *rose.Gold) func(string, string) {
				return gold.TOMLEq
			},
		},
		"yaml": {
			goldenFile: "yaml_eq.golden.yaml",
			input: `
jobs:
 test:
   runs-on: ubuntu-22.04
   steps:
     - uses: actions/checkout@v3
     - uses: actions/setup-go@v2
       with:
         go-version: "1.18"
`,
			test: func(g *rose.Gold) func(string, string) {
				return g.YAMLEq
			},
		},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			gold := rose.New(
				t,
				rose.WithPrefix("testdata", t.Name()),
				rose.WithFlag(*update),
			)
			f := tc.test(gold)
			f(tc.goldenFile, tc.input)
		})
	}
}
