package list

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/internal/view"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

// NewCmdList is a list command.
func NewCmdList() *cobra.Command {
	cmd := cobra.Command{
		Use:     "list",
		Short:   "List lists Jira projects",
		Long:    "List lists Jira projects that a user has access to.",
		Aliases: []string{"lists", "ls"},
		Run:     List,
	}
	cmd.Flags().Bool("raw", false, "Print the projects as JSON")
	return &cmd
}

// List displays a list view.
func List(cmd *cobra.Command, _ []string) {
	debug, err := cmd.Flags().GetBool("debug")
	cmdutil.ExitIfError(err)
	raw, err := cmd.Flags().GetBool("raw")
	cmdutil.ExitIfError(err)

	if raw {
		projects, err := api.DefaultClient(debug).Project()
		cmdutil.ExitIfError(err)
		if projects == nil {
			projects = []*jira.Project{}
		}
		out, err := json.Marshal(projects)
		cmdutil.ExitIfError(err)
		fmt.Println(string(out))
		return
	}

	projects, total, err := func() ([]*jira.Project, int, error) {
		s := cmdutil.Info("Fetching projects...")
		defer s.Stop()

		projects, err := api.DefaultClient(debug).Project()
		if err != nil {
			return nil, 0, err
		}
		return projects, len(projects), nil
	}()
	cmdutil.ExitIfError(err)

	if total == 0 {
		cmdutil.Failed("No projects found.")
		return
	}

	v := view.NewProject(projects)

	cmdutil.ExitIfError(v.Render())
}
