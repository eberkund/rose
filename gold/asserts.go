package gold

import (
	"github.com/eberkund/rose/formatting"
)

// JSONEq compares XML to golden file.
func (g *Gold) JSONEq(goldenPath, actual string) {
	g.t.Helper()
	g.verify(g.assert(goldenPath, actual, formatting.JSON))
}

// HTMLEq compares XML to golden file.
func (g *Gold) HTMLEq(goldenPath, actual string) {
	g.t.Helper()
	g.verify(g.assert(goldenPath, actual, formatting.HTML))
}

// TOMLEq compares TOML to golden file.
func (g *Gold) TOMLEq(goldenPath, actual string) {
	g.t.Helper()
	g.verify(g.assert(goldenPath, actual, formatting.TOML))
}

// YAMLEq compares YAML to golden file.
func (g *Gold) YAMLEq(goldenPath, actual string) {
	g.t.Helper()
	g.verify(g.assert(goldenPath, actual, formatting.YAML))
}

// Eq compares string to golden file.
func (g *Gold) Eq(goldenPath, actual string) {
	g.t.Helper()
	g.verify(g.assert(goldenPath, actual, formatting.NoOp))
}
