package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli"
	"google.golang.org/genproto/googleapis/devtools/cloudtrace/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"gtrace/span"
	"gtrace/tracer"
)

func projects(c *cli.Context) []string {
	var results []string
	for _, p := range c.StringSlice("project") {
		results = append(results, strings.Split(p, ",")...)
	}
	for _, p := range c.GlobalStringSlice("project") {
		results = append(results, strings.Split(p, ",")...)
	}
	return results
}

func get(projects []string, id string, writer io.Writer) error {
	ctx := context.Background()
	trc, err := tracer.NewTracer(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = trc.Close() }()

	trace, err := trc.MultiGet(ctx, projects, id)
	if err != nil {
		return err
	}

	span.Sort(trace.Spans)

	traceJSON, err := protojson.MarshalOptions{Indent: "\t"}.Marshal(trace)
	if err != nil {
		return err
	}

	_, err = writer.Write(traceJSON)
	if err != nil {
		return err
	}

	return nil
}

func list(project string, writer io.Writer, opt ...tracer.ListOption) error {
	ctx := context.Background()
	trc, err := tracer.NewTracer(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = trc.Close() }()

	fmt.Printf("listing...\n")
	traces, err := trc.List(ctx, project, opt...)
	if err != nil {
		return err
	}
	fmt.Printf("marshaling...\n")
	traceJSON, err := json.MarshalIndent(traces, "", "\t")
	if err != nil {
		return err
	}

	_, err = writer.Write(traceJSON)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	app := &cli.App{
		Name:      "gtrace",
		Version:   "v1.0.0",
		Compiled:  time.Now(),
		Copyright: "(c) 1999 Serious Enterprise",
		HelpName:  "contrive",
		Usage:     "demonstrate available API",
		UsageText: "contrive - demonstrating the available API",
		ArgsUsage: "[args and such]",
		Commands: []cli.Command{
			{
				Name: "get",
				Action: func(c *cli.Context) error {
					id := c.Args().First()
					proj := projects(c)
					return get(proj, id, os.Stdout)
				},
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name: "project",
					},
				},
			},
			{
				Name: "list",
				Action: func(c *cli.Context) error {
					proj := projects(c)
					if len(proj) <= 0 {
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

					req := &cloudtrace.ListTracesRequest{}
					fmt.Printf("project: %s\n", proj[0])
					for _, o := range opts {
						o(req)
					}
					fmt.Printf("%+v\n", req)
					return list(proj[0], os.Stdout, opts...)
				},
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name: "project",
					},
					&cli.IntFlag{
						Name:  "limit",
						Value: 10,
					},
					&cli.DurationFlag{
						Name: "since",
					},
					&cli.StringSliceFlag{
						Name: "filter",
					},
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name: "project",
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
