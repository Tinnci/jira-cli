// Package count counts issues for one or more JQL queries.
package count

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
)

// NewCmdCount is a count command.
func NewCmdCount() *cobra.Command {
	cmd := cobra.Command{
		Use:     "count",
		Short:   "Count prints how many issues match each JQL query",
		Long:    "Count prints a JSON array of {jql, total}, one per --jql, without reading any issue.",
		Example: `$ jira issue count --jql "project = PLV AND affectedVersion = V12" --jql "project = PLV AND affectedVersion = V14"`,
		Run:     count,
	}
	cmd.Flags().StringArray("jql", nil, "JQL query; repeat for several")
	return &cmd
}

type result struct {
	JQL   string `json:"jql"`
	Total int    `json:"total"`
}

func count(cmd *cobra.Command, _ []string) {
	debug, err := cmd.Flags().GetBool("debug")
	cmdutil.ExitIfError(err)
	queries, _ := cmd.Flags().GetStringArray("jql")
	if len(queries) == 0 {
		cmdutil.Failed("at least one --jql is required")
		os.Exit(1)
	}
	client := api.DefaultClient(debug)
	out := make([]result, 0, len(queries))
	for _, q := range queries {
		total, err := client.CountV2(q)
		cmdutil.ExitIfError(err)
		out = append(out, result{JQL: q, Total: total})
	}
	b, err := json.Marshal(out)
	cmdutil.ExitIfError(err)
	fmt.Println(string(b))
}
