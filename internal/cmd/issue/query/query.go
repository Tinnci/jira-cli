// Package query prints raw Jira search results.
package query

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
)

const maxLimit = 1000

// NewCmdQuery is a query command.
func NewCmdQuery() *cobra.Command {
	cmd := cobra.Command{
		Use:   "query",
		Short: "Query prints the REST search response for a JQL query",
		Long: `Query prints Jira's v2 /search response unchanged: total, startAt and each issue
with the requested fields, custom fields included. The JQL is used as given; --project is
not added to it.`,
		Example: `$ jira issue query --jql "project = PLV AND issuetype = Bug" --fields summary,status,customfield_15121 --limit 500`,
		Run:     query,
	}
	cmd.Flags().String("jql", "", "JQL query (required)")
	cmd.Flags().Uint("from", 0, "Index of the first issue")
	cmd.Flags().Uint("limit", 100, "Issues per call, at most 1000 (Jira may cap it lower)")
	cmd.Flags().String("fields", "", "Comma-separated fields; empty for Jira's default navigable fields")
	return &cmd
}

func query(cmd *cobra.Command, _ []string) {
	debug, err := cmd.Flags().GetBool("debug")
	cmdutil.ExitIfError(err)
	jql, _ := cmd.Flags().GetString("jql")
	from, _ := cmd.Flags().GetUint("from")
	limit, _ := cmd.Flags().GetUint("limit")
	fields, _ := cmd.Flags().GetString("fields")
	if jql == "" {
		cmdutil.Failed("--jql is required")
		os.Exit(1)
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	body, err := api.DefaultClient(debug).SearchRawV2(jql, from, limit, fields)
	cmdutil.ExitIfError(err)
	fmt.Println(string(body))
}
