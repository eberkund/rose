package rose

import (
	"testing"
)

func TestGold_JSONEq(t *testing.T) {
	g := New(t)
	g.JSONEq("testdata/json_eq.golden.json", `{"foo":123,"bar":"hello world","a":true}`)
}
