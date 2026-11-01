// Package filter lists saved Jira filters.
package filter

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

// NewCmdFilter is a filter command.
func NewCmdFilter() *cobra.Command {
	cmd := cobra.Command{
		Use:         "filter",
		Short:       "Filter lists saved Jira filters",
		Long:        "Filter lists saved Jira filters. See available commands below.",
		Annotations: map[string]string{"cmd:main": "true"},
		Aliases:     []string{"filters"},
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	list := cobra.Command{
		Use:     "list",
		Short:   "List the filters you starred, with their JQL",
		Aliases: []string{"ls"},
		Run:     list,
	}
	list.Flags().Bool("raw", false, "Print the filters as JSON")
	cmd.AddCommand(&list)
	return &cmd
}

func list(cmd *cobra.Command, _ []string) {
	debug, err := cmd.Flags().GetBool("debug")
	cmdutil.ExitIfError(err)
	raw, err := cmd.Flags().GetBool("raw")
	cmdutil.ExitIfError(err)

	filters, err := api.DefaultClient(debug).FavouriteFilters()
	cmdutil.ExitIfError(err)
	if filters == nil {
		filters = []*jira.Filter{}
	}
	if raw {
		out, err := json.Marshal(filters)
		cmdutil.ExitIfError(err)
		fmt.Println(string(out))
		return
	}
	if len(filters) == 0 {
		cmdutil.Failed("No starred filters found.")
		return
	}
	for _, f := range filters {
		fmt.Printf("%s\t%s\t%s\n", f.ID, f.Name, f.JQL)
	}
}
