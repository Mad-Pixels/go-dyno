package generator

import (
	"github.com/rs/zerolog"
)

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

	dryRun bool
}

// MarshalZerologObject return Config fields for logger.
func (c *Config) MarshalZerologObject(e *zerolog.Event) {
	e.Str("schemaPath", c.schemaPath)
	e.Str("outputPath", c.outputDir)
	e.Str("scope", string(c.scope))
	e.Bool("isDryRun", c.dryRun)

	if c.packageName != nil {
		e.Str("customPackageName", *c.packageName)
	}
	if c.fileName != nil {
		e.Str("customFileName", *c.fileName)
	}
}
