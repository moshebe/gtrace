package cli

import (
	"fmt"

	"github.com/moshebe/gtrace/pkg/span"
	"github.com/urfave/cli/v2"
	"google.golang.org/genproto/googleapis/devtools/cloudtrace/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

var formatAction = func(c *cli.Context) error {
	format := c.String("template")
	input, output := c.String("input"), c.String("output")

	in, err := read(input)
	if err != nil {
		return err
	}

	var trace cloudtrace.Trace
	err = protojson.Unmarshal(in, &trace)
	if err != nil {
		return fmt.Errorf("nmarshal trace: %w", err)
	}

	out, err := writer(output)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	return span.Format(trace.Spans, format, out)
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
		&cli.PathFlag{
			Name:    "output",
			Aliases: []string{"o", "out"},
			Value:   "-",
			Usage:   "output file path. '-' means stdout",
		},
		&cli.StringFlag{
			Name:  "template",
			Value: "{{ .Name }}  ({{ .Start }}  -  took {{ .Duration }})\n{{ if .Labels }}\t{{ .Labels }}\n{{ end }}",
			Usage: "templated pattern to format each span record base on TraceSpan properties\n\t",
		},
	},
}
