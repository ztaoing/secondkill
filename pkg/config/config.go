/**
* @Author:zhoutao
* @Date:2020/7/1 上午11:13
* @分布式配置中心接入组件
 */

package config

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"secondkill/pkg/bootstrap"
	"secondkill/pkg/discover"
	"strconv"
)

const KConfigType = "CONFIG_TYPE"

var Logger log.Logger
var zipkinTracer *zipkin.Tracer

func init() {
	Logger := log.NewLogfmtLogger(os.Stderr)
	Logger = log.With(Logger, "timestamp", log.DefaultTimestampUTC)
	Logger = log.With(Logger, "caller", log.DefaultCaller)

	viper.AutomaticEnv()
	initDeault()
	//加载远端配置
	if err := LoadRemoteConfig(); err != nil {
		Logger.Log("fail to load remote config", err)
	}
	if err := Sub("trace", &TraceConfig); err != nil {
		Logger.Log("Fail to parse trace", err)
	}

	zipkinUrl := "http://" + TraceConfig.Host + ":" + TraceConfig.Port + TraceConfig.Url
	Logger.Log("zipkin url", zipkinUrl)

	//初始化zipkin配置
	initTracer(zipkinUrl)
}

//设置配置文件类型
func initDeault() {
	viper.SetDefault(KConfigType, "yaml")
}

//加载远端配置
func LoadRemoteConfig() (err error) {
	//配置中心：服务发现+负载均衡
	serviceInstance, err := discover.DiscoveryService(bootstrap.ConfigServerConfig.ID)
	if err != nil {
		return
	}

	configServer := "http://" + serviceInstance.Host + ":" + strconv.Itoa(serviceInstance.Port)
	confAddr := fmt.Sprintf("%v/%v/%v-%v.%v",
		configServer,
		bootstrap.ConfigServerConfig.Label,
		bootstrap.DiscoverConfig.ServiceName,
		bootstrap.ConfigServerConfig.Profile,
		viper.Get(KConfigType),
	)
	resp, err := http.Get(confAddr)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//设置配置类型
	viper.SetConfigType(viper.GetString(KConfigType))
	if err := viper.ReadConfig(resp.Body); err != nil {
		return
	}
	Logger.Log("Load config from:", confAddr)
	return
}

//解析配置
func Sub(key string, value interface{}) error {
	Logger.Log("配置文件的前缀为：", key)
	sub := viper.Sub(key)
	sub.AutomaticEnv()
	sub.SetEnvPrefix(key)
	return sub.Unmarshal(value)
}

func initTracer(zipkinUrl string) {
	var (
		err           error
		useNoopTracer = (zipkinUrl == "")
		reporter      = zipkinhttp.NewReporter(zipkinUrl)
	)

	defer reporter.Close()

	zEndpoint, _ := zipkin.NewEndpoint(bootstrap.DiscoverConfig.ServiceName, bootstrap.DiscoverConfig.Port)
	zipkinTracer, err = zipkin.NewTracer(
		reporter, zipkin.WithLocalEndpoint(zEndpoint), zipkin.WithNoopTracer(useNoopTracer),
	)
	if err != nil {
		Logger.Log("new Tracer failed:", err)
		os.Exit(1)
	}
	if !useNoopTracer {
		Logger.Log("tracer", "zipkin", "type", "native", "url", zipkinUrl)
	}

}
