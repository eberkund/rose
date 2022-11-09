package gold_test

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
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
		extension string
		testFn    func(g *gold.Gold) func(string, string)
	}{
		"json": {
			extension: "json",
			testFn: func(g *gold.Gold) func(string, string) {
				return g.AssertEqualsJSON
			},
		},
		"text": {
			extension: "txt",
			testFn: func(g *gold.Gold) func(string, string) {
				return g.AssertEquals
			},
		},
		"xml": {
			extension: "xml",
			testFn: func(g *gold.Gold) func(string, string) {
				return g.AssertEqualsHTML
			},
		},
		"toml": {
			extension: "toml",
			testFn: func(g *gold.Gold) func(string, string) {
				return g.AssertEqualsTOML
			},
		},
		"yaml": {
			extension: "yaml",
			testFn: func(g *gold.Gold) func(string, string) {
				return g.AssertEqualsYAML
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(st *testing.T) {
			prefix := path.Join("testdata", t.Name(), name)
			g := gold.New(
				st,
				gold.WithPrefix(prefix),
				gold.WithFlag(*update),
				gold.WithFailAfter(),
			)
			testFn := tc.testFn(g)
			inputFile := path.Join(prefix, fmt.Sprintf("input.%s", tc.extension))
			goldenFile := fmt.Sprintf("golden.%s", tc.extension)
			data, _ := os.ReadFile(inputFile)
			testFn(goldenFile, string(data))
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
