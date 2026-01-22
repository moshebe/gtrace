package cli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/moshebe/gtrace/pkg/span"
	"github.com/moshebe/gtrace/pkg/tracer"
	"github.com/urfave/cli/v2"
	cloudtrace "cloud.google.com/go/trace/apiv1/tracepb"
)

type listResult struct {
	Span   string   `json:"name"`
	Traces []string `json:"traces"`
}

var listAction = func(c *cli.Context) error {
	if !c.IsSet("project") {
		return fmt.Errorf("missing project")
	}
	var limit int32 = 10

	opts := []tracer.ListOption{tracer.WithOnlyRootSpanView(), tracer.WithLimit(10)}

	if c.IsSet("limit") {
		limit = int32(c.Int("limit"))
		opts = append(opts, tracer.WithLimit(limit))
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

	traces, err := trc.List(ctx, c.String("project"), limit, opts...)
	if err != nil {
		return fmt.Errorf("list traces: %w", err)
	}

	rootSpans := span.ListRootSpans(traces)

	results := make([]listResult, 0, len(rootSpans))
	for name, ids := range rootSpans {
		results = append(results, listResult{
			Span:   name,
			Traces: ids,
		})
	}

	format := c.String("format")
	switch format {
	case "json":
		var output []byte
		if c.Bool("pretty") {
			output, err = json.MarshalIndent(results, "", "\t")
		} else {
			output, err = json.Marshal(results)
		}
		if err != nil {
			return fmt.Errorf("marshal results: %w", err)
		}
		fmt.Println(string(output))
	case "text":
		for _, result := range results {
			fmt.Printf("%s (%d traces)\n", result.Span, len(result.Traces))
			for _, traceID := range result.Traces {
				fmt.Printf("  %s\n", traceID)
			}
		}
	default:
		return fmt.Errorf("unsupported format: %s (supported formats: json, text)", format)
	}

	return nil
}

var ListCommand = &cli.Command{
	Name:      "list",
	Action:    listAction,
	Usage:     "Query traces from a project according to the given conditions",
	UsageText: "gtrace list [command options]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "project",
			Aliases: []string{"p"},
			Usage:   "the Google Cloud project ID to use for this invocation",
		},
		&cli.IntFlag{
			Name:  "limit",
			Value: 10,
			Usage: "maximum number of traces to return",
		},
		&cli.DurationFlag{
			Name:  "since",
			Usage: "time duration to inspect since now",
		},
		&cli.StringSliceFlag{
			Name:    "filter",
			Aliases: []string{"f"},
			Usage:   "filter traces according to Cloud Trace API syntax. can be set multiple times. See: https://cloud.google.com/trace/docs/trace-filters#filter_syntax",
		},
		&cli.TimestampFlag{
			Name:   "start",
			Layout: "2006-01-02T15:04:05",
			Usage:  "start of the time interval (inclusive) during which the trace data was collected from the application",
		},
		&cli.TimestampFlag{
			Name:   "end",
			Layout: "2006-01-02T15:04:05",
			Usage:  "end of the time interval (inclusive) during which the trace data was collected from the application",
		},
		&cli.BoolFlag{
			Name:  "pretty",
			Usage: "prettify JSON output",
		},
		&cli.StringFlag{
			Name:  "format",
			Value: "json",
			Usage: "output format: json or text",
		},
	},
}
