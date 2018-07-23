package outlyer

import (
	"fmt"
	"os"
	"os/user"

	"github.com/spf13/viper"
)

// UserConfig stores user-specific configurations
var UserConfig = viper.New()

func init() {
	user, err := user.Current()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not read user's home directory", err)
		os.Exit(1)
	}

	UserConfig.AddConfigPath(user.HomeDir)
	UserConfig.SetConfigName(".outlyer")
	UserConfig.SetDefault("headers.common.accept", "application/yaml")
	UserConfig.SetDefault("headers.common.user-agent", "outlyer/1.0")
	UserConfig.SetDefault("headers.post.content-type", "application/yaml")
	UserConfig.SetDefault("api-url", "https://api2.outlyer.com/v2")
	UserConfig.ReadInConfig()
}
