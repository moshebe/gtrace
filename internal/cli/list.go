package cli

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/moshebe/gtrace/pkg/span"
	"github.com/moshebe/gtrace/pkg/tracer"
	"github.com/urfave/cli/v2"
	"google.golang.org/genproto/googleapis/devtools/cloudtrace/v1"
)

var listAction = func(c *cli.Context) error {
	if !c.IsSet("project") {
		return fmt.Errorf("missing project")
	}

	opts := []tracer.ListOption{tracer.WithOnlyRootSpanView(), tracer.WithLimit(10)}

	if c.IsSet("limit") {
		opts = append(opts, tracer.WithLimit(int32(c.Int("limit"))))
	}

	if c.IsSet("since") {
		opts = append(opts, tracer.WithSince(c.Duration("since")))
	}

	if c.IsSet("filter") {
		opts = append(opts, tracer.WithFilter(c.StringSlice("filter")...))
	}

	if ts := c.Timestamp("start"); ts != nil {
		opts = append(opts, tracer.WithStartTime(*ts))
	}

	if ts := c.Timestamp("end"); ts != nil {
		opts = append(opts, tracer.WithEndTime(*ts))
	}

	req := &cloudtrace.ListTracesRequest{}
	for _, o := range opts {
		o(req)
	}

	ctx := context.Background()
	trc, err := tracer.NewTracer(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = trc.Close() }()

	traces, err := trc.List(ctx, c.String("project"), opts...)
	if err != nil {
		return fmt.Errorf("list traces: %w", err)
	}

	rootSpans := span.ListRootSpans(traces)

	w := tabwriter.NewWriter(os.Stdout, 0, 150, 10, '\t', tabwriter.AlignRight)
	for name, ids := range rootSpans {
		if len(ids) > 1 {
			_, err = fmt.Fprintf(w, "%s\t\t[%v + %d more...]\n", name, ids[0], len(ids)-1)
		} else {
			_, err = fmt.Fprintf(w, "%s\t\t%v\n", name, ids)
		}
		if err != nil {
			return err
		}
	}
	_ = w.Flush()

	return nil
}

var ListCommand = &cli.Command{
	Name:   "list",
	Action: listAction,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "project",
			Aliases: []string{"p"},
		},
		&cli.IntFlag{
			Name:  "limit",
			Value: 10,
		},
		&cli.DurationFlag{
			Name: "since",
		},
		&cli.StringSliceFlag{
			Name:    "filter",
			Aliases: []string{"f"},
		},
		&cli.TimestampFlag{
			Name:   "start",
			Layout: "2006-01-02T15:04:05",
		},
		&cli.TimestampFlag{
			Name:   "end",
			Layout: "2006-01-02T15:04:05",
		},
		&cli.BoolFlag{
			Name: "pretty",
		},
	},
}
