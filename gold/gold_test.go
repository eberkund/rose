package gold_test

import (
	"flag"
	"io"
	"testing"

	"github.com/eberkund/rose/gold"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "update golden files")

func TestExistingFiles(t *testing.T) {
	testCases := map[string]struct {
		inputData  string
		goldenFile string
		testFn     func(g *gold.Gold) func(string, string)
	}{
		"json": {
			inputData:  `{"foo":123,"bar":"hello world","a":true}`,
			goldenFile: "json_eq.golden.json",
			testFn: func(gold *gold.Gold) func(string, string) {
				return gold.JSONEq
			},
		},
		"text": {
			goldenFile: "text_eq.golden.txt",
			inputData:  "Hello\nWorld\n!",
			testFn: func(gold *gold.Gold) func(string, string) {
				return gold.Eq
			},
		},
		"html": {
			goldenFile: "xml_eq.golden.toml",
			inputData:  `<fruits><apple/><banana/></fruits>`,
			testFn: func(gold *gold.Gold) func(string, string) {
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
			testFn: func(gold *gold.Gold) func(string, string) {
				return gold.TOMLEq
			},
		},
		"yaml": {
			goldenFile: "yaml_eq.golden.yaml",
			inputData: `
jobs:
 test:
   runs-on: ubuntu-22.04
   steps:
     - uses: actions/checkout@v3
     - uses: actions/setup-go@v2
       with:
         go-version: "1.18"
`,
			testFn: func(g *gold.Gold) func(string, string) {
				return g.YAMLEq
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gold := gold.New(
				t,
				gold.WithPrefix("testdata", t.Name()),
				gold.WithFlag(*update),
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

	g := gold.New(t, gold.WithFS(memFs), gold.WithFlag(true), gold.WithFatal(false))
	g.Eq("test_data.txt", newData)

	data, err := afero.ReadFile(memFs, "testdata/test_data.txt")
	require.NoError(t, err)
	require.Equal(t, newData, string(data))
}

//func TestFailingDiff(t *testing.T) {
//	gold := rose.New(t, rose.WithPrefix("testdata", "TestExistingFiles", "json"))
//	gold.JSONEq("json_eq.golden.json", "{}")
//	if t.Failed() {
//		t.SkipNow()
//		println("very nice!")
//	}
//}
