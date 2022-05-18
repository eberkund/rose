# Rose 

Rose is a golden file utility library for your Go tests.

- Excellent support for JSON, XML, YAML


## Usage

```go
package mypackage

import (
	"testing"
)

func (t *testing.T) {
	reader := someReader()
	rose.Golden(t, "somefile.txt", reader)
	rose.GoldenJSON(t, "somefile.json", reader)
	rose.GoldenXML(t, "", reader)
	rose.GoldenYAML(t, "", reader)
	rose.GoldenTOML(t, "", reader)
}
```

or with config...

```go
package mypackage

import (
	"testing"
)

func (t *testing.T) {
    g := rose.NewGolden(
        t, 
        rose.PrettyFormat(),
        rose.IgnoreOrder(),
    )
    g.JSON("somefile.json", reader)
}
```

## Available Methods

- JSONEqual
- XMLEqual
- TOMLEqual
- YAMLEqual
- Equal