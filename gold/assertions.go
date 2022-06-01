package gold

import (
	"github.com/eberkund/rose/formatting"
)

// AssertEqualsJSON compares a JSON string to the golden file. Data is
// noramlized before comparing so data which is formatted differently but
// semantically equivalent will not fail.
func (g *Gold) AssertEqualsJSON(goldenPath, actual string) {
	g.t.Helper()
	g.verify(g.assert(goldenPath, actual, formatting.JSON))
}

// AssertEqualsHTML compares an HTML string to the golden file. Data is
// noramlized before comparing so data which is formatted differently but
// semantically equivalent will not fail.
func (g *Gold) AssertEqualsHTML(goldenPath, actual string) {
	g.t.Helper()
	g.verify(g.assert(goldenPath, actual, formatting.HTML))
}

// AssertEqualsTOML compares a TOML string to the golden file. Data is
// noramlized before comparing so data which is formatted differently but
// semantically equivalent will not fail.
func (g *Gold) AssertEqualsTOML(goldenPath, actual string) {
	g.t.Helper()
	g.verify(g.assert(goldenPath, actual, formatting.TOML))
}

// AssertEqualsYAML compares a YAML string to the golden file. Data is
// noramlized before comparing so data which is formatted differently but
// semantically equivalent will not fail.
func (g *Gold) AssertEqualsYAML(goldenPath, actual string) {
	g.t.Helper()
	g.verify(g.assert(goldenPath, actual, formatting.YAML))
}

// AssertEquals compares a string to the contents of the golden file. If Gold
// was configured with an update flag which is true then the golden file will
// be created if it does not exist or overwritten with the supplied data if it
// does.
func (g *Gold) AssertEquals(goldenPath, actual string) {
	g.t.Helper()
	g.verify(g.assert(goldenPath, actual, formatting.NoOp))
}
