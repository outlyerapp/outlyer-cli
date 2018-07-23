package command

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"

	"github.com/outlyer/outlyer-cli"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

// NewConfigureCommand creates a Command for setting up the user's local
// Outlyer yaml configuration file given the API token
func NewConfigureCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "configure",
		Short: "Set up the Outlyer CLI by validating your API token",
		Run:   createLocalConfig,
	}
	return cmd
}

type cliConfig struct {
	APIToken string `yaml:"api-token"`
}

// createLocalConfig validates the API token provided by the user
// and persists it locally by creating the a hidden Outlyer yaml configuration file
// at the user's $HOME directory.
func createLocalConfig(cmd *cobra.Command, args []string) {
	for i := 0; i < 3; i++ {
		// Reads the API token from stdin
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Please enter your API token: ")
		apiToken, _ := reader.ReadString('\n')
		apiToken = strings.Replace(apiToken, "\n", "", -1) // removes return character on *unix and darwin
		apiToken = strings.Replace(apiToken, "\r", "", -1) // removes return character on windows

		outlyer.UserConfig.Set("api-token", apiToken)

		_, err := outlyer.Get("/user")
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Errorf("Error: invalid API token\n%s", err))
			continue
		}

		userConfig := cliConfig{apiToken}

		userConfigInBytes, err := yaml.Marshal(userConfig)
		if err != nil {
			ExitWithError(ExitError, err)
		}

		user, err := user.Current()
		if err != nil {
			ExitWithError(ExitError, err)
		}

		err = ioutil.WriteFile(user.HomeDir+"/.outlyer.yaml", userConfigInBytes, 0644)
		if err != nil {
			log.Fatalln(err)
		}
		ExitWithSuccess("Success! Outlyer CLI is configured and ready to use")
	}
	ExitWithError(ExitError, fmt.Errorf("Could not configure Outlyer CLI. Please contact Outlyer support"))
}
