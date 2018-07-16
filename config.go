package outlyer

import (
	"os"

	"github.com/spf13/viper"
)

// UserConfig stores user-specific configurations
var UserConfig = viper.New()

func init() {
	UserConfig.AddConfigPath(os.Getenv("HOME"))
	UserConfig.SetConfigName(".outlyer")
	UserConfig.SetDefault("headers.common.accept", "application/text")
	UserConfig.SetDefault("headers.common.user-agent", "outlyer/1.0")
	UserConfig.SetDefault("headers.post.content-type", "application/yaml")
	UserConfig.SetDefault("api-url", "https://api2.outlyer.com/v2")
	UserConfig.ReadInConfig()
}
