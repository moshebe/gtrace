package cli

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/moshebe/gtrace/pkg/span"
	"github.com/urfave/cli/v2"
	cloudtrace "cloud.google.com/go/trace/apiv1/tracepb"
	"google.golang.org/protobuf/encoding/protojson"
)

var durationAction = func(c *cli.Context) error {
	file := c.String("file")
	min := c.Duration("min")

	if min == time.Duration(0) {
		return fmt.Errorf("missing minimum duration")
	}

	in, err := read(file)
	if err != nil {
		return err
	}

	var trace cloudtrace.Trace
	err = protojson.Unmarshal(in, &trace)
	if err != nil {
		return fmt.Errorf("unmarshal trace: %w", err)
	}

	trace.Spans = span.FilterMinDuration(trace.Spans, min)

	if c.Bool("sort") {
		sort.Slice(trace.Spans, func(i, j int) bool {
			return span.Duration(trace.Spans[i]) > span.Duration(trace.Spans[j])
		})
	}

	if !c.Bool("summary") {
		return printTraceJSON(os.Stdout, &trace)
	}

	for _, s := range trace.Spans {
		fmt.Println(span.DurationSummary(s))
	}
	return nil
}

var DurationCommand = &cli.Command{
	Name:  "duration",
	Usage: "Filter trace spans by total duration",
	Description: "The duration is calculated by subtracting start from the end time. Can be useful to troubleshoot " +
		"performance or latency issues within a specific trace",
	UsageText: "gtrace duration [command options]",
	Action:    durationAction,
	Flags: []cli.Flag{
		&cli.PathFlag{
			Name:    "file",
			Aliases: []string{"f"},
			Value:   "-",
			Usage:   "input file path. '-' means stdin",
		},
		&cli.DurationFlag{
			Name:  "min",
			Value: time.Second,
			Usage: "minimum duration threshold for filtering",
		},
		&cli.BoolFlag{
			Name:  "summary",
			Value: true,
			Usage: "output a short summary of the results",
		},
		&cli.BoolFlag{
			Name:  "sort",
			Value: true,
			Usage: "sort the results in descending order by span duration",
		},
	},
}
