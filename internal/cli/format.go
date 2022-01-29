package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/moshebe/gtrace/pkg/span"
	"github.com/urfave/cli/v2"
	"google.golang.org/genproto/googleapis/devtools/cloudtrace/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

var formatAction = func(c *cli.Context) error {
	format := c.String("template")
	input, output := c.String("input"), c.String("output")

	in, out := os.Stdin, os.Stdout

	if !stdio(input) {
		f, err := os.Open(input)
		if err != nil {
			return err
		}
		defer func() { _ = f.Close() }()

		in = f
	}

	if !stdio(output) {
		f, err := os.OpenFile(output, createFileFlags, createFilePerm)
		if err != nil {
			return err
		}
		defer func() { _ = f.Close() }()

		out = f
	}

	var trace cloudtrace.Trace
	traceJSON, err := io.ReadAll(in)
	if err != nil {
		return fmt.Errorf("read input: %w", err)
	}
	err = protojson.Unmarshal(traceJSON, &trace)
	if err != nil {
		return fmt.Errorf("nmarshal trace: %w", err)
	}

	return span.Format(trace.Spans, format, out)
}

var FormatCommand = &cli.Command{
	Name:   "format",
	Action: formatAction,
	Flags: []cli.Flag{
		&cli.PathFlag{
			Name:    "input",
			Aliases: []string{"i", "in"},
			Value:   "-",
		},
		&cli.PathFlag{
			Name:    "output",
			Aliases: []string{"o", "out"},
			Value:   "-",
		},
		&cli.StringFlag{
			Name:  "template",
			Value: "{{ .Name }}  ({{ .Start }}  -  took {{ .Duration }})\n{{ if .Labels }}\t{{ .Labels }}\n{{ end }}",
		},
	},
}
