package flags

import (
	"fmt"
	"strings"

	godyno "github.com/Mad-Pixels/go-dyno"

	"github.com/urfave/cli/v2"
)

var (
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
)
