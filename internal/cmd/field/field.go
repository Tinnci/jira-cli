// Package field lists the Jira instance's fields.
package field

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

// NewCmdField is a field command.
func NewCmdField() *cobra.Command {
	cmd := cobra.Command{
		Use:         "field",
		Short:       "Field lists the fields configured in Jira",
		Long:        "Field lists the system and custom fields configured in Jira, with their IDs.",
		Annotations: map[string]string{"cmd:main": "true"},
		Aliases:     []string{"fields"},
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	list := cobra.Command{
		Use:     "list",
		Short:   "List field IDs and names, e.g. customfield_17302 → 问题模块",
		Aliases: []string{"ls"},
		Run:     list,
	}
	list.Flags().Bool("raw", false, "Print the fields as JSON")
	cmd.AddCommand(&list)
	return &cmd
}

func list(cmd *cobra.Command, _ []string) {
	debug, err := cmd.Flags().GetBool("debug")
	cmdutil.ExitIfError(err)
	raw, err := cmd.Flags().GetBool("raw")
	cmdutil.ExitIfError(err)

	fields, err := api.DefaultClient(debug).GetField()
	cmdutil.ExitIfError(err)
	if fields == nil {
		fields = []*jira.Field{}
	}
	if raw {
		out, err := json.Marshal(fields)
		cmdutil.ExitIfError(err)
		fmt.Println(string(out))
		return
	}
	for _, f := range fields {
		fmt.Printf("%s\t%s\n", f.ID, f.Name)
	}
}
