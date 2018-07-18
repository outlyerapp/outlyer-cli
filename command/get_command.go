package command

import "github.com/spf13/cobra"

// NewGetCommand groups subcommands like "accounts"
func NewGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get the specified resources. The available resources are: 'accounts'",
	}
	cmd.AddCommand(NewGetAccountCommand())
	return cmd
}
