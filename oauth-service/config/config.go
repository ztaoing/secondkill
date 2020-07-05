/**
* @Author:zhoutao
* @Date:2020/7/2 下午9:27
加载配置
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

var (
	ZipKinTracer *zipkin.Tracer
	Logger       log.Logger
)

const KConfigType = "CONFIG_TYPE"

func init() {
	//1.日志
	Logger = log.NewLogfmtLogger(os.Stderr)
	Logger = log.With(Logger, "timestamp", log.DefaultTimestampUTC)
	Logger = log.With(Logger, "caller", log.DefaultCaller)

	viper.AutomaticEnv()
	initDefault()

	//2.加载远端配置
	if err := config.LoadRemoteConfig(); err != nil {
		Logger.Log("fail to load remote config from remote service", err)
	}

	if err := config.Sub("mysql", &config.MysqlConfig); err != nil {
		Logger.Log("fail to parse mysql ", err)
	}

	if err := config.Sub("trace", &config.TraceConfig); err != nil {
		Logger.Log("fail to parse trace ", err)
	}

	//3.zipkin链路追踪
	zipkinUrl := "http://" + config.TraceConfig.Host + ":" + config.TraceConfig.Port + config.TraceConfig.Url
	Logger.Log("zipkin url", zipkinUrl)
	initTracer(zipkinUrl)
}

func initDefault() {
	viper.SetDefault(KConfigType, "yaml")
}

func initTracer(zipkinUrl string) {
	var (
		err           error
		useNoopTracer = (zipkinUrl == "")
		reporter      = zipkinhttp.NewReporter(zipkinUrl)
	)
	defer reporter.Close()

	zEndpoint, _ := zipkin.NewEndpoint(bootstrap.DiscoverConfig.ServiceName, bootstrap.DiscoverConfig.Port)
	_, err = zipkin.NewTracer(
		reporter, zipkin.WithLocalEndpoint(zEndpoint), zipkin.WithNoopTracer(useNoopTracer),
	)
	if err != nil {
		Logger.Log("new tracer failed", err)
		//0正常退出，1非正常退出
		os.Exit(1)
	}
	if !useNoopTracer {
		Logger.Log("tracer", "zipkin", "type", "native", "url", zipkinUrl)
	}
}
