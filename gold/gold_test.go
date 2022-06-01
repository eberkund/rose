package gold_test

import (
	"flag"
	"io"
	"testing"

	"github.com/eberkund/rose/gold"
	"github.com/eberkund/rose/mocks"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/mock"
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
			testFn: func(g *gold.Gold) func(string, string) {
				return g.JSONEq
			},
		},
		"text": {
			goldenFile: "text_eq.golden.txt",
			inputData:  "Hello\nWorld\n!",
			testFn: func(g *gold.Gold) func(string, string) {
				return g.Eq
			},
		},
		"html": {
			goldenFile: "xml_eq.golden.toml",
			inputData:  `<fruits><apple/><banana/></fruits>`,
			testFn: func(g *gold.Gold) func(string, string) {
				return g.HTMLEq
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
			testFn: func(g *gold.Gold) func(string, string) {
				return g.TOMLEq
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
			g := gold.New(
				t,
				gold.WithPrefix("testdata", t.Name()),
				gold.WithFlag(*update),
			)
			f := tc.testFn(g)
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

func TestFailingDiff(t *testing.T) {
	tm := mocks.NewTesting(t)
	tm.On("Helper")
	tm.On("Logf", mock.Anything, mock.Anything).Once()
	tm.On("FailNow").Once()
	g := gold.New(tm, gold.WithPrefix("testdata", "TestExistingFiles", "json"))
	g.JSONEq("json_eq.golden.json", "{}")
}

func TestFilesystemFails(t *testing.T) {
	tm := mocks.NewTesting(t)
	tm.On("Helper")
	tm.On("Log", mock.Anything).Once()
	tm.On("FailNow").Once()
	readOnlyFS := afero.NewReadOnlyFs(afero.NewMemMapFs())
	g := gold.New(tm, gold.WithFS(readOnlyFS), gold.WithFlag(true))
	g.Eq("testdata/hello_world.txt", "Hello, World!")
}
