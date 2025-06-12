package generator

type scope string

const (
	// ScopeAll for build full pkg.
	ScopeAll scope = "all"
)

// Config holds finalized generator configuration.
type Config struct {
	schemaPath string
	outputDir  string
	scope      scope

	packageName *string // optional: override inferred package name
	fileName    *string // optional: override inferred file name

	verbose bool
	dryRun  bool
}
