package main

import (
	"context"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli"
	"google.golang.org/protobuf/encoding/protojson"
	"gtrace/span"
	"gtrace/trace"
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

// TODO: move from here
func get(projects []string, id string, writer io.Writer) error {
	ctx := context.Background()
	tracer, err := trace.NewTracer(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tracer.Close() }()

	trc, err := tracer.MultiGet(ctx, projects, id)
	if err != nil {
		return err
	}

	span.Sort(trc.Spans)

	traceJSON, err := protojson.MarshalOptions{Indent: "\t"}.Marshal(trc)
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
