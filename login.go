package outlyer

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/spf13/viper"

	yaml "gopkg.in/yaml.v2"
)

type cliConfig struct {
	APIToken       string `yaml:"api-token"`
	DefaultAccount string `yaml:"default-account"`
}

// CreateLocalConfig uses the API token to fetch the default user account
// and creates a hidden Outlyer yaml configuration file locally in the user's $HOME directory.
func CreateLocalConfig() error {
	resp, err := Get("/user")
	if err != nil {
		return err
	}

	viper.SetConfigType("yaml")
	err = viper.ReadConfig(bytes.NewReader(resp))
	if err != nil {
		return err
	}

	defaultAccount := viper.GetString("default_account")
	userConfig := cliConfig{UserConfig.GetString("api-token"), defaultAccount}

	userConfigYaml, err := yaml.Marshal(userConfig)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(os.Getenv("HOME")+"/.outlyer.yaml", userConfigYaml, 0644)
	if err != nil {
		return err
	}

	return nil
}
