/**
* @Author:zhoutao
* @Date:2020/7/2 下午2:57
统一业务api网关
*/

package main

import (
	"flag"
	"github.com/go-kit/kit/log"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"os"
	"secondkill/pkg/bootstrap"
	register "secondkill/pkg/discover"
)

func main() {
	var zipkinURL = flag.String("zipkin.url", "http://127.0.0.1:9411/api/v2/spans", "zipkin server url")
	flag.Parse()

	//创建日志组件
	var Logger log.Logger
	{
		Logger = log.NewLogfmtLogger(os.Stderr)
		Logger = log.With(Logger, "timestamp", log.DefaultTimestampUTC)
		Logger = log.With(Logger, "caller", log.DefaultCaller)
	}

	//设置zipkin组件
	var zipkinTracer *zipkin.Tracer
	{
		var (
			err           error
			useNoopTracer = (*zipkinURL == "")
			reporter      = zipkinhttp.NewReporter(*zipkinURL)
		)
		defer reporter.Close()

		zEndpoint, _ := zipkin.NewEndpoint(bootstrap.HttpConfig.Host, bootstrap.DiscoverConfig.Port)
		zipkinTracer, err = zipkin.NewTracer(
			reporter, zipkin.WithLocalEndpoint(zEndpoint), zipkin.WithNoopTracer(useNoopTracer),
		)

		if err != nil {
			Logger.Log("err", err.Error())
		}

		if !useNoopTracer {
			Logger.Log("tracer", "zipkin", "type", "native", "url", *zipkinURL)
		}
	}

	//注册服务
	register.Register()

}
