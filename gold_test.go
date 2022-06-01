package rose_test

import (
	"flag"
	"io"
	"testing"

	"github.com/eberkund/rose"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "update golden files")

func TestExistingFiles(t *testing.T) {
	testcases := map[string]struct {
		input      string
		goldenFile string
		test       func(g *rose.Gold) func(string, string, ...interface{})
	}{
		"json": {
			input:      `{"foo":123,"bar":"hello world","a":true}`,
			goldenFile: "json_eq.golden.json",
			test: func(gold *rose.Gold) func(string, string, ...interface{}) {
				return gold.JSONEq
			},
		},
		"text": {
			goldenFile: "text_eq.golden.txt",
			input:      "Hello\nWorld\n!",
			test: func(gold *rose.Gold) func(string, string, ...interface{}) {
				return gold.Eq
			},
		},
		"html": {
			goldenFile: "xml_eq.golden.toml",
			input:      `<fruits><apple/><banana/></fruits>`,
			test: func(gold *rose.Gold) func(string, string, ...interface{}) {
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
			test: func(gold *rose.Gold) func(string, string, ...interface{}) {
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
			test: func(g *rose.Gold) func(string, string, ...interface{}) {
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

func TestUpdate(t *testing.T) {
	const (
		oldData = "foo"
		newData = "bar"
	)

	memFs := afero.NewMemMapFs()
	file, err := memFs.Create("testdata/test_data.txt")
	require.NoError(t, err)
	_, err = io.WriteString(file, oldData)
	require.NoError(t, err)
	err = file.Close()
	require.NoError(t, err)

	g := rose.New(t, rose.WithFS(memFs), rose.WithFlag(true))
	g.Eq("test_data.txt", newData)

	data, err := afero.ReadFile(memFs, "testdata/test_data.txt")
	require.NoError(t, err)
	require.Equal(t, newData, string(data))
}

//func TestFailingDiff(t *testing.T) {
//	gold := rose.New(t, rose.WithPrefix("testdata", "TestExistingFiles", "json"))
//	gold.JSONEq("json_eq.golden.json", "{}")
//}
