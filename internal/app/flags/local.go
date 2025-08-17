package flags

import (
	"fmt"
	"strings"

	godyno "github.com/Mad-Pixels/go-dyno"
	"github.com/Mad-Pixels/go-dyno/internal/generator/mode"

	"github.com/urfave/cli/v2"
)

var (
	// LocalSchema defines the --schema flag for specifying the input JSON schema file.
	// This flag is required for all commands that need to process a DynamoDB schema definition.
	LocalSchema = Flag{
		Object: &cli.StringFlag{
			Name:  "schema",
			Usage: "Set path to 'JSON' schame for generate goLang objects.",
			Aliases: []string{
				"s",
			},
			EnvVars: []string{
				fmt.Sprintf("%s_%s", godyno.EnvPrefix, strings.ToUpper("schema")),
			},
			Required: true,
		},
	}

	// LocalOutputDir defines the --output-dir flag for specifying where generated files should be written.
	LocalOutputDir = Flag{
		Object: &cli.StringFlag{
			Name:  "output-dir",
			Usage: "Set destination directory path. (write to stdout if not set)",
			Aliases: []string{
				"o",
			},
			EnvVars: []string{
				fmt.Sprintf("%s_%s", godyno.EnvPrefix, strings.ToUpper("output-dir")),
			},
			Required: false,
		},
	}

	// LocalPackageName defines the --package flag for overriding the generated Go package name.
	// By default, the package name is derived from the table_name in the schema.
	LocalPackageName = Flag{
		Object: &cli.StringFlag{
			Name:    "package",
			Usage:   "Overwrite generated file package name. (default is 'table_name' value)",
			Aliases: []string{},
			EnvVars: []string{
				fmt.Sprintf("%s_%s", godyno.EnvPrefix, strings.ToUpper("package")),
			},
			Required: false,
		},
	}

	// LocalFilename defines the --filename flag for overriding the generated Go file name.
	// By default, the filename is derived from the table_name in the schema with .go extension.
	LocalFilename = Flag{
		Object: &cli.StringFlag{
			Name:    "filename",
			Usage:   "Overwrite generated filename. (default is 'table_name' value)",
			Aliases: []string{},
			EnvVars: []string{
				fmt.Sprintf("%s_%s", godyno.EnvPrefix, strings.ToUpper("filename")),
			},
			Required: false,
		},
	}

	// LocalGenerateMode defines the --mode flag for controlling code generation mode.
	// Determines what code to generate: ALL (complete) or MIN (minimal).
	LocalGenerateMode = Flag{
		Object: &cli.StringFlag{
			Name:  "mode",
			Usage: fmt.Sprintf("Set generation mode (%s). (default: %s)", strings.Join(mode.GetAvailableModes(), ", "), mode.GetDefault()),
			Aliases: []string{
				"m",
			},
			EnvVars: []string{
				fmt.Sprintf("%s_%s", godyno.EnvPrefix, strings.ToUpper("mode")),
			},
			Required: false,
			Value:    mode.GetDefault().String(),
		},
	}

	// LocalWithStreamEvents defines the --with-stream-events for methods which work with DynamoDB stream
	// By default, stream events methods not included.
	LocalWithStreamEvents = Flag{
		Object: &cli.BoolFlag{
			Name:    "with-stream-events",
			Usage:   "Add methods with works with DynamoDB streams",
			Aliases: []string{},
			EnvVars: []string{
				fmt.Sprintf("%s_%s", godyno.EnvPrefix, strings.ToUpper("with-stream-events")),
			},
			Required: false,
		},
	}
)
