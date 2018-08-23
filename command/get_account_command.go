package command

import (
	"fmt"

	"github.com/outlyerapp/outlyer-cli/api"
	"github.com/spf13/cobra"
)

// NewGetAccountCommand creates a Command to list the user's accounts
func NewGetAccountCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accounts",
		Short: "List user accounts",
		Run:   listUserAccounts,
	}
	return cmd
}

// listUserAccounts fetches the Outlyer API and lists the user's accounts
func listUserAccounts(cmd *cobra.Command, args []string) {
	resp, err := api.Get("/accounts")
	if err != nil {
		ExitWithError(ExitError, fmt.Errorf("Could not fetch user's accounts\n%s", err))
	}
	ExitWithSuccess(string(resp))
}
