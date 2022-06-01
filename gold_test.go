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
	testCases := map[string]struct {
		inputData  string
		goldenFile string
		testFn     func(g *rose.Gold) func(string, string, ...interface{})
	}{
		"json": {
			inputData:  `{"foo":123,"bar":"hello world","a":true}`,
			goldenFile: "json_eq.golden.json",
			testFn: func(gold *rose.Gold) func(string, string, ...interface{}) {
				return gold.JSONEq
			},
		},
		"text": {
			goldenFile: "text_eq.golden.txt",
			inputData:  "Hello\nWorld\n!",
			testFn: func(gold *rose.Gold) func(string, string, ...interface{}) {
				return gold.Eq
			},
		},
		"html": {
			goldenFile: "xml_eq.golden.toml",
			inputData:  `<fruits><apple/><banana/></fruits>`,
			testFn: func(gold *rose.Gold) func(string, string, ...interface{}) {
				return gold.HTMLEq
			},
		},
		"toml": {
			goldenFile: "toml_eq.golden.toml",
			inputData: `
Age = 25
Cats = [ "Cauchy", "Plato" ]

Pi = 3.14
Perfection = [ 6, 28, 496, 8128 ]
DOB = 1987-07-05T05:45:00Z
`,
			testFn: func(gold *rose.Gold) func(string, string, ...interface{}) {
				return gold.TOMLEq
			},
		},
		"yaml": {
			goldenFile: "yaml_eq.golden.yaml",
			inputData: `
jobs:
 testFn:
   runs-on: ubuntu-22.04
   steps:
     - uses: actions/checkout@v3
     - uses: actions/setup-go@v2
       with:
         go-version: "1.18"
`,
			testFn: func(g *rose.Gold) func(string, string, ...interface{}) {
				return g.YAMLEq
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gold := rose.New(
				t,
				rose.WithPrefix("testdata", t.Name()),
				rose.WithFlag(*update),
			)
			f := tc.testFn(gold)
			f(tc.goldenFile, tc.inputData)
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
