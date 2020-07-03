/**
* @Author:zhoutao
* @Date:2020/7/2 下午2:57
统一业务api网关
*/

package main

import (
	"flag"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-kit/kit/log"
	"github.com/openzipkin/zipkin-go"
	zipkinhttpSvr "github.com/openzipkin/zipkin-go/middleware/http"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"net"
	"net/http"
	"os"
	"os/signal"
	"secondkill/gateway/route"
	"secondkill/pkg/bootstrap"
	register "secondkill/pkg/discover"
	"syscall"
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

	tags := map[string]string{
		"component": "gateway_server",
	}

	hystrixRouter := route.Routes(zipkinTracer, "Circuit Breaker:service unavalable", Logger)

	handler := zipkinhttpSvr.NewServerMiddleware(
		zipkinTracer,
		zipkinhttpSvr.SpanName(bootstrap.DiscoverConfig.ServiceName),
		zipkinhttpSvr.TagResponseSize(true),
		zipkinhttpSvr.ServerTags(tags),
	)(hystrixRouter)

	errc := make(chan error)

	//启用hystrix实时监控，监听端口为9010
	hystrixSreamHandler := hystrix.NewStreamHandler()
	hystrixSreamHandler.Start()

	go func() {
		errc <- http.ListenAndServe(net.JoinHostPort("", "9010"), hystrixSreamHandler)

	}()

	//监听停止信号
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	//开始监听
	go func() {
		Logger.Log("stransport", "HTTP", "addr", "9090")
		register.Register()
		errc <- http.ListenAndServe(":9090", handler)
	}()
}
