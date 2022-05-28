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
			goldenFile: "assert_json.golden.json",
			testFn: func(g *gold.Gold) func(string, string) {
				return g.AssertEqualsJSON
			},
		},
		"text": {
			goldenFile: "assert_text.golden.txt",
			inputData:  "Hello\nWorld\n!",
			testFn: func(g *gold.Gold) func(string, string) {
				return g.AssertEquals
			},
		},
		"html": {
			goldenFile: "assert_xml.golden.toml",
			inputData:  `<fruits><apple/><banana/></fruits>`,
			testFn: func(g *gold.Gold) func(string, string) {
				return g.AssertEqualsHTML
			},
		},
		"toml": {
			goldenFile: "assert_toml.golden.toml",
			inputData: `
Age = 25
Cats = [ "Cauchy", "Plato" ]

Pi = 3.14
Perfection = [ 6, 28, 496, 8128 ]
DOB = 1987-07-05T05:45:00Z
`,
			testFn: func(g *gold.Gold) func(string, string) {
				return g.AssertEqualsTOML
			},
		},
		"yaml": {
			goldenFile: "assert_yaml.golden.yaml",
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
				return g.AssertEqualsYAML
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(st *testing.T) {
			g := gold.New(
				st,
				gold.WithPrefix("testdata", t.Name()),
				gold.WithFlag(*update),
				gold.WithFailAfter(),
			)
			f := tc.testFn(g)
			f(tc.goldenFile, tc.inputData)
		})
	}
}

func TestUpdatingGoldenFile(t *testing.T) {
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

	g := gold.New(t, gold.WithFS(memFs), gold.WithFlag(true))
	g.AssertEquals("test_data.txt", newData)

	data, err := afero.ReadFile(memFs, "testdata/test_data.txt")
	require.NoError(t, err)
	require.Equal(t, newData, string(data))
}

func TestInvalidData(t *testing.T) {
	tm := mocks.NewTesting(t)
	tm.On("Helper")
	tm.On("Log", mock.Anything).Once()
	tm.On("FailNow").Once()
	g := gold.New(tm, gold.WithPrefix("testdata", t.Name()))
	g.AssertEqualsJSON("empty_json.golden.json", "this is not valid JSON")
}

func TestDiffGoldenFile(t *testing.T) {
	tm := mocks.NewTesting(t)
	tm.On("Helper")
	tm.On("Logf", mock.Anything, mock.Anything).Once()
	tm.On("FailNow").Once()
	g := gold.New(tm, gold.WithPrefix("testdata", t.Name()))
	g.AssertEqualsJSON("empty_json.golden.json", "{}")
}

func TestFileSystemError(t *testing.T) {
	tm := mocks.NewTesting(t)
	tm.On("Helper")
	tm.On("Log", mock.Anything).Once()
	tm.On("FailNow").Once()
	readOnlyFS := afero.NewReadOnlyFs(afero.NewMemMapFs())
	g := gold.New(tm, gold.WithFS(readOnlyFS), gold.WithFlag(true))
	g.AssertEquals("any/path/hello_world.golden.txt", "Hello, World!")
}
