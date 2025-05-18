package generate

import (
	"bytes"
	"strings"
	"text/template"

	godyno "github.com/Mad-Pixels/go-dyno"
	cli "github.com/urfave/cli/v2"
)

var (
	name  = "gen"
	usage = "generate static golang code from config"
)

type tmplUsage struct {
}

func Command() *cli.Command {
	tmpl, err := template.New("usage").Funcs(template.FuncMap{
		"Join": strings.Join,
	}).Parse(usageTemplate)
	if err != nil {
		godyno.Logger.Fatal(err)
	}

	var bText bytes.Buffer
	err = tmpl.Execute(&bText, tmplUsage{})
	if err != nil {
		godyno.Logger.Fatal(err)
	}

	return &cli.Command{
		Name:      name,
		Usage:     usage,
		UsageText: bText.String(),
		Action:    action,
		Flags:     flags(),
	}
}
