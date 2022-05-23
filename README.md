[![Tests](https://github.com/eberkund/rose/actions/Tests/badge.svg)](https://github.com/eberkund/rose/actions/workflows/Tests)
[![codecov](https://codecov.io/gh/eberkund/rose/branch/master/graph/badge.svg?token=lCcKXaBzlD)](https://codecov.io/gh/eberkund/rose)
[![Go Report Card](https://goreportcard.com/badge/github.com/eberkund/rose)](https://goreportcard.com/report/github.com/eberkund/rose)

# Rose

Rose is a golden file utility library for your Go tests.

It handles the dirty work so your tests stay concise and readable.

- Normalizes common file formats
- Helps organize files with per-test path prefixes
- Leaves you in control of supplying the update flag

## Usage

```go
package mypackage

import (
	"flag"
	"testing"

	"github.com/eberkund/rose/gold"
)

// Use `go test ./... -update` to write new files with supplied data
var update = flag.Bool("update", false, "update golden files")

func TestFiles(t *testing.T) {
    g := gold.New(
        t,

        // Pass the update flag to Gold constructor
        gold.WithFlag(*update),

        // Store files in `testdata/<test name>/<file name>`
        // Prefix defaults to `testdata`
        gold.WithPrefix("testdata", t.Name()),
    )

    // Provide the filename and input data
    g.JSONEq("somefile.golden.json", `{"foo": 123}`)
}
```
