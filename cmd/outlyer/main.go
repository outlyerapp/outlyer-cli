package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	outlyer "github.com/outlyer/outlyer-cli"
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "outlyer",
	Short: "Outlyer CLI allows to export and apply your Outlyer account via command line.",
	Run: func(cmd *cobra.Command, args []string) {
		log.Fatalln("Must use a subcommand")
	},
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Sets up the Outlyer CLI given the API token.",
	Run: func(cmd *cobra.Command, args []string) {
		for i := 0; i < 3; i++ {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Please enter your API token: ")
			apiToken, _ := reader.ReadString('\n')
			apiToken = strings.Replace(apiToken, "\n", "", -1)
			outlyer.UserConfig.Set("api-token", apiToken)
			err := outlyer.CreateLocalConfig()
			if err != nil {
				fmt.Printf("Error validating API Token!\n%s\n\n", err)
			} else {
				fmt.Println("\nSuccess! Outlyer CLI is configured and ready to use.")
				return
			}
		}
		fmt.Println("Please contact Outlyer support.")
		os.Exit(1)
	},
}

func main() {
	cmd.AddCommand(loginCmd)
	cmd.Execute()
}
