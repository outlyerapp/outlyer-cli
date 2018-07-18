package command

import (
	"fmt"
	"log"

	"github.com/outlyer/outlyer-cli"
	"github.com/spf13/cobra"
)

// NewGetAccountCommand creates a Command to list user accounts
func NewGetAccountCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accounts",
		Short: "List user accounts",
		Run:   listUserAccounts,
	}
	return cmd
}

// listUserAccounts lists the user's accounts
func listUserAccounts(cmd *cobra.Command, args []string) {
	resp, err := outlyer.Get("/accounts")
	if err != nil {
		log.Fatalln("Error fetching user accounts", err)
	}
	fmt.Println(string(resp))
}
