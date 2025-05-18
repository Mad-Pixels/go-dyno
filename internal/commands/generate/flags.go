package generate

import (
	"fmt"

	godyno "github.com/Mad-Pixels/go-dyno"
	"github.com/urfave/cli/v2"
)

func flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "qwe",
			Usage: "usage",
			EnvVars: []string{
				fmt.Sprintf("%s_AA", godyno.EnvPrefix),
			},
		},
	}
}
