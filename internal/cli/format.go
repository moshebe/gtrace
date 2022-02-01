package cli

import (
	"fmt"
	"os"

	"github.com/moshebe/gtrace/pkg/span"
	"github.com/urfave/cli/v2"
	"google.golang.org/genproto/googleapis/devtools/cloudtrace/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

var formatAction = func(c *cli.Context) error {
	format := c.String("template")
	input := c.String("input")

	in, err := read(input)
	if err != nil {
		return err
	}

	var trace cloudtrace.Trace
	err = protojson.Unmarshal(in, &trace)
	if err != nil {
		return fmt.Errorf("unmarshal trace: %w", err)
	}

	return span.Format(trace.Spans, format, os.Stdout)
}

var FormatCommand = &cli.Command{
	Name:        "format",
	Usage:       "Format trace spans according to a given template",
	Description: "See more information at: https://cloud.google.com/trace/docs/reference/v1/rest/v1/projects.traces#TraceSpan",
	UsageText:   "gtrace format [command options]",
	Action:      formatAction,
	Flags: []cli.Flag{
		&cli.PathFlag{
			Name:    "input",
			Aliases: []string{"i", "in"},
			Value:   "-",
			Usage:   "input file path. '-' means stdin",
		},
		&cli.StringFlag{
			Name:  "template",
			Value: "{{ .Name }}  ({{ .Start }}  -  took {{ .Duration }})\n{{ if .Labels }}\t{{ .Labels }}\n{{ end }}",
			Usage: "templated pattern to format each span record base on TraceSpan properties\n\t",
		},
	},
}
