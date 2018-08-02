package config

import (
	"fmt"
	"os"
	"os/user"

	"github.com/spf13/viper"
)

// CLI stores Outlyer configurations
var CLI = viper.New()

func init() {
	user, err := user.Current()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not read user's home directory", err)
		os.Exit(1)
	}

	CLI.AddConfigPath(user.HomeDir)
	CLI.SetConfigName(".outlyer")
	CLI.SetDefault("headers.common.accept", "application/yaml")
	CLI.SetDefault("headers.common.user-agent", "outlyer/1.0")
	CLI.SetDefault("headers.post.content-type", "application/yaml")
	CLI.SetDefault("api-url", "https://api2.outlyer.com/v2")
	CLI.ReadInConfig()
}
