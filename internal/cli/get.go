package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/moshebe/gtrace/pkg/span"
	"github.com/moshebe/gtrace/pkg/tracer"
	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/encoding/protojson"
)

var getAction = func(c *cli.Context) error {
	id := c.Args().First()
	output := c.Path("output")
	projects := stringSlice(c, "project")

	if len(projects) == 0 {
		return fmt.Errorf("missing project")
	}
	if id == "" {
		return fmt.Errorf("missing trace id")
	}

	out, err := writer(output)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	ctx, cancel := context.WithTimeout(c.Context, time.Minute)
	defer cancel()

	trc, err := tracer.NewTracer(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = trc.Close() }()

	trace, err := trc.MultiGet(ctx, projects, id)
	if err != nil {
		return fmt.Errorf("get trace: %w", err)
	}

	span.Sort(trace.Spans)

	indent := ""
	if c.Bool("pretty") {
		indent = "\t"
	}
	traceJSON, err := protojson.MarshalOptions{Indent: indent}.Marshal(trace)
	if err != nil {
		return fmt.Errorf("marshal trace: %w", err)
	}

	_, err = out.Write(traceJSON)
	if err != nil {
		return fmt.Errorf("write trace: %w", err)
	}

	return nil
}

var GetCommand = &cli.Command{
	Name:        "get",
	Action:      getAction,
	Usage:       "Get a specific trace by id from one or more projects",
	Description: "Retrieve the trace information from the given project(s), aggregate the results and sort the spans by their start time.",
	UsageText:   "gtrace get [command options] <trace-id>",
	Flags: []cli.Flag{
		&cli.StringSliceFlag{
			Name:    "project",
			Aliases: []string{"p"},
			Usage:   "the Google Cloud project ID to use for this invocation. values can be set multiple times or separated by comma",
		},
		&cli.PathFlag{
			Name:    "output",
			Aliases: []string{"o", "out"},
			Value:   "-",
			Usage:   "output file path. '-' means stdout",
		},
		&cli.BoolFlag{
			Name:  "pretty",
			Usage: "prettify JSON output",
		},
	},
}
