/**
* @Author:zhoutao
* @Date:2020/7/6 下午3:08
 */

package config

import (
	"github.com/go-kit/kit/log"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/spf13/viper"
	"os"
	"secondkill/pkg/bootstrap"
	"secondkill/pkg/config"
)

var ZipkinTracer *zipkin.Tracer
var Logger log.Logger

const KconfigType = "CONFIG_TYPE"

func init() {
	Logger := log.NewLogfmtLogger(os.Stderr)
	Logger = log.With(Logger, "timestamp", log.DefaultTimestamp)
	Logger = log.With(Logger, "caller", log.DefaultCaller)

	viper.AutomaticEnv()
	initDefault()

	//加载远端的配置
	if err := config.LoadRemoteConfig(); err != nil {
		Logger.Log("faild to load remote config", err)
	}
	//从config中解析
	if err := config.Sub("mysql", &config.MysqlConfig); err != nil {
		Logger.Log("failed to parse mysql from config", err)
	}
	if err := config.Sub("trace", &config.TraceConfig); err != nil {
		Logger.Log("failed to parse trace from config", err)
	}
	if err := config.Sub("redis", &config.Redis); err != nil {
		Logger.Log("failed to parse redis from config", err)
	}
	if err := config.Sub("service", &config.SecKill); err != nil {
		Logger.Log("failed to parse service from config", err)
	}

	zipkinUrl := "http://" + config.TraceConfig.Host + ":" + config.TraceConfig.Port + config.TraceConfig.Url

	Logger.Log("zipkin url", zipkinUrl)
	initTracer(zipkinUrl)

}

func initDefault() {
	viper.SetDefault(KconfigType, "yaml")
}

func initTracer(zipkinUrl string) {
	var (
		err           error
		useNoopTracer = (zipkinUrl == "")
		reporter      = zipkinhttp.NewReporter(zipkinUrl)
	)
	defer reporter.Close()

	//endpoint
	zEndpoint, err := zipkin.NewEndpoint(bootstrap.DiscoverConfig.ServiceName, bootstrap.HttpConfig.Port)
	ZipkinTracer, err = zipkin.NewTracer(reporter,
		zipkin.WithLocalEndpoint(zEndpoint),
		zipkin.WithNoopTracer(useNoopTracer),
	)
	if err != nil {
		Logger.Log("err", err)
		os.Exit(1)
	}
	if !useNoopTracer {
		Logger.Log("tracer", "zipkin", "type", "native", "url", zipkinUrl)
	}

}
