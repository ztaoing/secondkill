/**
* @Author:zhoutao
* @Date:2020/7/1 上午7:49
 */

package bootstrap

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

func init() {
	viper.AutomaticEnv()
	//初始化读取配置文件、路径、配置类型
	initBootstrapConfig()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("err:%s\n", err)
	}
	//解析http配置
	if err := subParse("http", &HttpConfig); err != nil {
		log.Fatal("fail to parse http config", err)
	}
	//解析discover配置
	if err := subParse("discover", &DiscoverConfig); err != nil {
		log.Fatal("fail to parse discover config", err)
	}
	//解析配置中心配置
	if err := subParse("config", &ConfigServerConfig); err != nil {
		log.Fatal("fail to parse config server", err)
	}
	//解析PRC配置
	if err := subParse("rpc", &RpcConfig); err != nil {
		log.Fatal("fail to parse rpc server", err)
	}
}

func subParse(key string, value interface{}) error {
	log.Printf("配置文件的前缀为：%v", key)
	sub := viper.Sub(key)
	sub.AutomaticEnv()
	sub.SetEnvPrefix(key)
	return sub.Unmarshal(value)
}

func initBootstrapConfig() {
	//设置读取的配置文件
	viper.SetConfigName("bootstrap")
	//添加读取配置文件的路径
	viper.AddConfigPath("./")
	//gopath路径下的配置
	viper.AddConfigPath("$GOPATH/src/")
	//设置配置文件类型
	viper.SetConfigType("yaml")
}
