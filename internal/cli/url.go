package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var urlAction = func(c *cli.Context) error {
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
}

var URLCommand = &cli.Command{
	Name:      "url",
	Usage:     "Generate a browsable URL for a given trace.",
	UsageText: "gtrace get [--project <project-id>] <trace-id>",
	Action:    urlAction,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "project",
			Aliases: []string{"p"},
			Usage:   "the Google Cloud project ID to specify on the generated URL.",
		},
	},
}
