package command

import (
	"fmt"
	"regexp"

	"github.com/spf13/cobra"
)

type plugin struct {
	Content  string `yaml:"content"`
	Encoding string `yaml:"encoding"`
	Name     string `yaml:"name"`
}

type resource struct {
	path   string
	bytes  []byte
	status string
	err    error
}

func (r *resource) getType() string {
	regex := regexp.MustCompile(`(alerts|checks|dashboards|plugins)`)
	res := regex.FindStringSubmatch(r.path)
	return res[0]
}

func (r *resource) getTypeAndName() string {
	regex := regexp.MustCompile(`(alerts|checks|dashboards|plugins)/[^.]+`)
	res := regex.FindStringSubmatch(r.path)
	return res[0]
}

func (r *resource) getTypeAndNameWithExtension() string {
	regex := regexp.MustCompile(`(alerts|checks|dashboards|plugins)/.+`)
	res := regex.FindStringSubmatch(r.path)
	return res[0]
}

func (r *resource) getNameWithExtension() string {
	regex := regexp.MustCompile(`(.*)/(.+)$`)
	res := regex.FindStringSubmatch(r.path)
	return res[2]
}

// NewApplyCommand creates a Command for applying resources to the user's Outlyer account
func NewApplyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply .|[folder]|[file]",
		Short: "Apply resources to the specified account. The available resources are: alerts, checks, dashboards and plugins",
		Run:   applyCommand,
	}

	cmd.PersistentFlags().StringP("account", "a", "", "(Required) User account to use")
	return cmd
}

func applyCommand(cmd *cobra.Command, args []string) {
	account := cmd.PersistentFlags().Lookup("account").Value.String()
	if account == "" {
		ExitWithError(ExitBadArgs, fmt.Errorf("Account is required"))
	}

	if len(args) < 1 {
		ExitWithError(ExitBadArgs, fmt.Errorf("Resource is required"))
	}
}
