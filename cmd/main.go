package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
	"google.golang.org/genproto/googleapis/devtools/cloudtrace/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"gtrace/span"
	"gtrace/tracer"
)

const (
	createFileFlags = os.O_RDWR | os.O_CREATE
	createFilePerm  = 0660
)

func stdio(value string) bool {
	return value == "-" || value == ""
}

func stringSlice(c *cli.Context, name string) []string {
	var results []string
	for _, v := range c.StringSlice(name) {
		results = append(results, strings.Split(v, ",")...)
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

	rootSpans := span.ListRootSpans(traces)
	for name, ids := range rootSpans {
		_, err := fmt.Fprintf(writer, "%s (%v)\n", name, ids)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	app := &cli.App{
		Name:      "gtrace",
		Version:   "v1.0.0",
		Compiled:  time.Now(),
		Copyright: "(c) 1999 Serious Enterprise",
		HelpName:  "gtrace",
		Usage:     "demonstrate available API",
		UsageText: "contrive - demonstrating the available API",
		ArgsUsage: "[args and such]",
		Commands: []*cli.Command{
			{
				Name: "get",
				Action: func(c *cli.Context) error {
					id := c.Args().First()
					output := c.Path("output")
					projects := stringSlice(c, "project")

					if id == "" {
						return fmt.Errorf("missing trace id")
					}

					writer := os.Stdout
					if output != "-" && output != "" {
						f, err := os.OpenFile(output, createFileFlags, createFilePerm)
						if err != nil {
							return err
						}
						defer func() { _ = f.Close() }()
						writer = f
					}

					return get(projects, id, writer)
				},
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
			},
			{
				Name: "list",
				Action: func(c *cli.Context) error {
					projects := stringSlice(c, "project")
					if len(projects) <= 0 {
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
					fmt.Printf("%+v\n", req)
					return list(projects[0], os.Stdout, opts...)
				},
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
			},
			{
				Name: "url",
				Action: func(c *cli.Context) error {
					id := c.Args().First()
					if id == "" {
						return fmt.Errorf("missing trace id")
					}

					projectPath := ""
					if c.IsSet("project") {
						projectPath += "&project=" + c.String("project")
					}
					fmt.Printf("https://console.cloud.google.com/traces/list?tid=%s%s\n", id, projectPath)
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "project",
						Aliases: []string{"p"},
					},
				},
			},
			{
				Name: "format",
				Action: func(c *cli.Context) error {
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
				},
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
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
