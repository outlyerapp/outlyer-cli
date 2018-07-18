package command

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/outlyer/outlyer-cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

// NewConfigureCommand validates the API token and persists it locally
// in the user's configuration file.
func NewConfigureCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "configure",
		Short: "Set up the Outlyer CLI by validating your API token",
		Run:   createLocalConfig,
	}
	return cmd
}

type cliConfig struct {
	APIToken       string `yaml:"api-token"`
	DefaultAccount string `yaml:"default-account"`
}

// createLocalConfig uses the API token provided by the user to fetch the default user account
// and creates a hidden Outlyer yaml configuration file in the user's $HOME directory.
func createLocalConfig(cmd *cobra.Command, args []string) {
	for i := 0; i < 3; i++ {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Please enter your API token: ")
		apiToken, _ := reader.ReadString('\n')
		apiToken = strings.Replace(apiToken, "\n", "", -1)
		outlyer.UserConfig.Set("api-token", apiToken)

		resp, err := outlyer.Get("/user")
		if err != nil {
			fmt.Println("Error validating the API token.", err)
			continue
		}

		viper.SetConfigType("yaml")
		err = viper.ReadConfig(bytes.NewReader(resp))
		if err != nil {
			log.Fatalln(err)
		}

		defaultAccount := viper.GetString("default_account")
		userConfig := cliConfig{apiToken, defaultAccount}

		userConfigYaml, err := yaml.Marshal(userConfig)
		if err != nil {
			log.Fatalln(err)
		}

		err = ioutil.WriteFile(os.Getenv("HOME")+"/.outlyer.yaml", userConfigYaml, 0644)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("\nSuccess! Outlyer CLI is configured and ready to use.\n")
		return
	}
	fmt.Println("Please contact Outlyer support.")
	os.Exit(1)
}
