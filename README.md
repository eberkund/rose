# Rose 

[![.github/workflows/test.yml](https://github.com/eberkund/rose/actions/workflows/test.yml/badge.svg)](https://github.com/eberkund/rose/actions/workflows/test.yml)

Rose is a golden file utility library for your Go tests.

- Excellent support for JSON, XML, YAML


## Usage

```go
package mypackage

import (
	"flag"
	"testing"
	
    "github.com/eberkund/rose"
)

var update bool

func init() {
	flag.BoolVar(&update, "update", false, "update golden files")
}

func TestFiles(t *testing.T) {
	gold := rose.New(t, rose.UpdateFlag(update))
	gold.Eq(t, "somefile.golden.txt", "hello\nworld\n!")
	gold.JSONEq(t, "somefile.golden.json", reader)
	gold.TOMLEq(t, "somefile.golden.toml", reader)
}

```


## Available Methods

- JSONEqual
- XMLEqual
- TOMLEqual
- YAMLEqual
- Equal
