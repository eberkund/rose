package rose

import (
	"strings"
	"testing"
)

func TestGold_JSONEq(t *testing.T) {
	g := New(t)
	g.JSONEq("testdata/json_eq.golden.json", strings.NewReader(`{"foo":123,"bar":"hello world"}`))
}
