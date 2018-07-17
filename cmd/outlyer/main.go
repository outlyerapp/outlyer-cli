package main

import (
	"github.com/outlyer/outlyer-cli/command"
	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command

func init() {
	rootCmd = &cobra.Command{
		Use:   "outlyer",
		Short: "Outlyer CLI allows to easily manage your Outlyer account via command line",
	}
	rootCmd.AddCommand(command.ConfigureCommand())
}

func main() {
	rootCmd.Execute()
}
