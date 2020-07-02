/**
* @Author:zhoutao
* @Date:2020/7/2 下午2:39
 */

package config

import (
	"github.com/go-kit/kit/log"
	"github.com/spf13/viper"
	"os"
	"secondkill/pkg/config"
)

var Logger log.Logger

const kConfigType = "CONFIG_TYPE"

func init() {
	Logger = log.NewLogfmtLogger(os.Stderr)
	Logger = log.With(Logger, "timestamp", log.DefaultTimestampUTC)
	Logger = log.With(Logger, "caller", log.DefaultCaller)

	viper.AutomaticEnv()
	initDefault()

	if err := config.LoadRemoteConfig(); err != nil {
		Logger.Log("fail to load remote config", err)
	}
	if err := config.Sub("auth", &AuthPermitConfig); err != nil {
		Logger.Log("fail to parse config", err)
	}
}

func initDefault() {
	viper.SetDefault(kConfigType, "yaml")
}
