/**
* @Author:zhoutao
* @Date:2020/7/7 上午8:09
 */

package config

import (
	"github.com/go-kit/kit/log"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/spf13/viper"
	"os"
	"secondkill/pkg/bootstrap"
	pkgConfig "secondkill/pkg/config"
)

const KConfigType = "CONFIG_TYPE"

var zipkinTracer *zipkin.Tracer
var Logger log.Logger

func init() {
	Logger = log.NewLogfmtLogger(os.Stderr)
	Logger = log.With(Logger, "timestamp", log.DefaultTimestampUTC)
	Logger = log.With(Logger, "caller", log.DefaultCaller)

	viper.AutomaticEnv()
	initDefault()

	//读取远端配置
	if err := pkgConfig.LoadRemoteConfig(); err != nil {
		Logger.Log("load remote config failed,ERROR:%v", err)
	}
	if err := pkgConfig.Sub("mysql", &pkgConfig.MysqlConfig); err != nil {
		Logger.Log("parse mysql failed,ERROR:%v", err)
	}
	if err := pkgConfig.Sub("trace", &pkgConfig.TraceConfig); err != nil {
		Logger.Log("parse trace failed,ERROR:%v", err)
	}

	zipkinURL := "http://" + pkgConfig.TraceConfig.Host + ":" + pkgConfig.TraceConfig.Port + pkgConfig.TraceConfig.Url
	Logger.Log("zipkin url", zipkinURL)

	initTracer(zipkinURL)
}

func initDefault() {
	viper.SetDefault(KConfigType, "yaml")
}

func initTracer(zipkinURL string) {
	var (
		err           error
		useNoopTracer = (zipkinURL == "")
		reporter      = zipkinhttp.NewReporter(zipkinURL)
	)
	defer reporter.Close()
	ZEndpoint, err := zipkin.NewEndpoint(bootstrap.DiscoverConfig.ServiceName, bootstrap.DiscoverConfig.Port)
	zipkinTracer, err = zipkin.NewTracer(reporter,
		zipkin.WithLocalEndpoint(ZEndpoint),
		zipkin.WithNoopTracer(useNoopTracer),
	)
	if err != nil {
		Logger.Log("err", err)
		os.Exit(1)
	}
	if !useNoopTracer {
		Logger.Log("tracer", "zipkin", "type", "native", "url", zipkinURL)
	}
}
