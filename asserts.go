package rose

// JSONEq compares XML to golden file.
func (g *Gold) JSONEq(goldenPath, actual string, msgAndArgs ...interface{}) {
	g.handleError(g.assert(goldenPath, actual, formatJSON))
}

// HTMLEq compares XML to golden file.
func (g *Gold) HTMLEq(goldenPath, actual string, msgAndArgs ...interface{}) {
	g.handleError(g.assert(goldenPath, actual, formatHTML))
}

// TOMLEq compares TOML to golden file.
func (g *Gold) TOMLEq(goldenPath, actual string, msgAndArgs ...interface{}) {
	g.handleError(g.assert(goldenPath, actual, formatTOML))
}

// YAMLEq compares YAML to golden file.
func (g *Gold) YAMLEq(goldenPath, actual string, msgAndArgs ...interface{}) {
	g.handleError(g.assert(goldenPath, actual, formatYAML))
}

// Eq compares string to golden file.
func (g *Gold) Eq(goldenPath, actual string, msgAndArgs ...interface{}) {
	g.handleError(g.assert(goldenPath, actual, formatNoop))
}
