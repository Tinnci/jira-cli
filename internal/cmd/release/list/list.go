package list

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/internal/view"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

// NewCmdList is a list command.
func NewCmdList() *cobra.Command {
	cmd := cobra.Command{
		Use:     "list",
		Short:   "List lists Jira projects versions",
		Long:    "List lists Jira projects versions that a user has access to.",
		Aliases: []string{"lists", "ls"},
		Run:     List,
	}
	cmd.Flags().Bool("raw", false, "Print the versions as JSON")
	return &cmd
}

// List displays a list view.
func List(cmd *cobra.Command, _ []string) {
	project := viper.GetString("project.key")
	debug, err := cmd.Flags().GetBool("debug")
	cmdutil.ExitIfError(err)
	raw, err := cmd.Flags().GetBool("raw")
	cmdutil.ExitIfError(err)

	if raw {
		releases, err := api.ProxyRelease(api.DefaultClient(debug), project)
		cmdutil.ExitIfError(err)
		if releases == nil {
			releases = []*jira.ProjectVersion{}
		}
		out, err := json.Marshal(releases)
		cmdutil.ExitIfError(err)
		fmt.Println(string(out))
		return
	}

	releases, total, err := func() ([]*jira.ProjectVersion, int, error) {
		s := cmdutil.Info("Fetching project versions...")
		defer s.Stop()

		releases, err := api.ProxyRelease(api.DefaultClient(debug), project)
		if err != nil {
			return nil, 0, err
		}
		return releases, len(releases), nil
	}()
	cmdutil.ExitIfError(err)

	if total == 0 {
		cmdutil.Failed("No releases found.")
		return
	}

	v := view.NewRelease(releases)

	cmdutil.ExitIfError(v.Render())
}
