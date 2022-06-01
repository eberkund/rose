package rose

// JSONEq compares XML to golden file.
func (g *Gold) JSONEq(goldenPath, actual string) {
	g.genericEQ(goldenPath, actual, formatJSON)
}

// XMLEq compares XML to golden file.
func (g *Gold) XMLEq(goldenPath, actual string) {
	g.genericEQ(goldenPath, actual, formatXML)
}

// TOMLEq compares TOML to golden file.
func (g *Gold) TOMLEq(goldenPath, actual string) {
	g.genericEQ(goldenPath, actual, formatTOML)
}

// YAMLEq compares YAML to golden file.
func (g *Gold) YAMLEq(goldenPath, actual string) {
	g.genericEQ(goldenPath, actual, formatYAML)
}

// Eq compares string to golden file.
func (g *Gold) Eq(goldenPath, actual string) {
	g.genericEQ(goldenPath, actual, formatNoop)
}
