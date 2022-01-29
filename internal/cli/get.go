package cli

import (
	"context"
	"fmt"
	"os"

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

	writer := os.Stdout
	if !stdio(output) {
		f, err := os.OpenFile(output, createFileFlags, createFilePerm)
		if err != nil {
			return fmt.Errorf("open file %q: %w", output, err)
		}
		defer func() { _ = f.Close() }()
		writer = f
	}

	ctx := context.Background()
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

	traceJSON, err := protojson.MarshalOptions{Indent: "\t"}.Marshal(trace)
	if err != nil {
		return fmt.Errorf("marshal trace: %w", err)
	}

	_, err = writer.Write(traceJSON)
	if err != nil {
		return fmt.Errorf("write trace: %w", err)
	}

	return nil
}

var GetCommand = &cli.Command{
	Name:   "get",
	Action: getAction,
	Flags: []cli.Flag{
		&cli.StringSliceFlag{
			Name:    "project",
			Aliases: []string{"p"},
		},
		&cli.PathFlag{
			Name:    "output",
			Aliases: []string{"o", "out"},
			Value:   "-",
		},
		&cli.BoolFlag{
			Name: "pretty",
		},
	},
}
