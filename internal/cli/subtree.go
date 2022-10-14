package cli

import (
	"fmt"
	"os"

	"github.com/moshebe/gtrace/pkg/span"
	"github.com/urfave/cli/v2"
	"google.golang.org/genproto/googleapis/devtools/cloudtrace/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

var subtreeAction = func(c *cli.Context) error {
	file := c.String("file")
	root := c.Uint64("root")

	if root == 0 {
		return fmt.Errorf("missing root span id")
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
	trace.Spans, err = span.SubTree(trace.Spans, root)
	if err != nil {
		return err
	}
	return printTraceJSON(os.Stdout, &trace)
}

var SubtreeCommand = &cli.Command{
	Name:  "subtree",
	Usage: "Extract span and all its children for a given trace",
	Description: "Iterating the trace spans according to the original ordering and opt-out all spans that " +
		"are not descendants the given root span",
	UsageText: "gtrace subtree [command options]",
	Action:    subtreeAction,
	Flags: []cli.Flag{
		&cli.PathFlag{
			Name:    "file",
			Aliases: []string{"f"},
			Value:   "-",
			Usage:   "input file path. '-' means stdin",
		},
		&cli.Uint64Flag{
			Name:  "root",
			Value: 0,
			Usage: "root span id",
		},
	},
}
