# Rose 

[![.github/workflows/test.yml](https://github.com/eberkund/rose/actions/workflows/test.yml/badge.svg)](https://github.com/eberkund/rose/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/eberkund/rose/branch/master/graph/badge.svg?token=lCcKXaBzlD)](https://codecov.io/gh/eberkund/rose)
[![Go Report Card](https://goreportcard.com/badge/github.com/eberkund/rose)](https://goreportcard.com/report/github.com/eberkund/rose)

Rose is a golden file utility library for your Go tests.

- Excellent support for JSON, XML, YAML

## Usage

```go
package mypackage

import (
	"flag"
	"strings"
	"testing"

	"github.com/eberkund/rose"
)

var update bool

func init() {
	flag.BoolVar(&update, "update", false, "update golden files")
}

func TestFiles(t *testing.T) {
	gold := rose.New(t, rose.UpdateFlag(update), rose.Prefix("testdata"))
	gold.Eq(t, "somefile.golden.txt", "hello\nworld\n!")
	gold.JSONEq(t, "somefile.golden.json", strings.NewReader(reader))
	gold.TOMLEq(t, "somefile.golden.toml", strings.NewReader(reader))
}

```


## Available Methods

- JSONEqual
- XMLEqual
- TOMLEqual
- YAMLEqual
- Equal
