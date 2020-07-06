/**
* @Author:zhoutao
* @Date:2020/7/6 上午7:42
 */

package setup

import (
	"context"
	"flag"
	"fmt"
	kitPrometheus "github.com/go-kit/kit/metrics/prometheus"
	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	stdPrometheus "github.com/prometheus/client_golang/prometheus"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"os"
	"os/signal"
	localConfig "secondkill/oauth-service/config"
	register "secondkill/pkg/discover"
	"secondkill/sk-app/config"
	"secondkill/sk-app/endpoint"
	"secondkill/sk-app/plugins"
	"secondkill/sk-app/service"
	"secondkill/sk-app/transport"
	"syscall"
	"time"
)

//初始化http服务
func InitHTTP(host, servicePort string) {
	log.Printf("host:%s port:%s", host, servicePort)

	flag.Parse()

	errChan := make(chan error)

	fiedlKeys := []string{"method"}

	requestCount := kitPrometheus.NewCounterFrom(stdPrometheus.CounterOpts{
		Namespace: "secKill-app",
		Subsystem: "sk-app",
		Name:      "request_count",
		Help:      "Number of requests received",
	}, fiedlKeys)

	//潜在的请求
	requestLatency := kitPrometheus.NewSummaryFrom(stdPrometheus.SummaryOpts{
		Namespace: "secKill-adpp",
		Subsystem: "sk_app",
		Name:      "request_latency",
		Help:      "Total duration of requests in microseconds",
	}, fiedlKeys)

	rateBucket := rate.NewLimiter(rate.Every(time.Second*1), 5000)

	var skAppService service.Service
	skAppService = service.SkAppService{}

	//日志组件
	skAppService = plugins.NewSkAppLoggingMiddleware(config.Logger)(skAppService)
	//metrics组件
	skAppService = plugins.NewSkAppMetricsMiddleware(requestCount, requestLatency)(skAppService)

	//healthCheckEndpoint

	healthCheckEndpoint := endpoint.MakeHealthCheckEndpoint(skAppService)
	healthCheckEndpoint = plugins.NewTokenBuketLimitterWithBuildIn(rateBucket)(healthCheckEndpoint)
	healthCheckEndpoint = kitzipkin.TraceEndpoint(localConfig.ZipKinTracer, "health-check-zipkin")(healthCheckEndpoint)

	//SecInfoEndpoint

	GetSecInfoEndpoint := endpoint.MakeSecInfoEndpoint(skAppService)
	GetSecInfoEndpoint = plugins.NewTokenBuketLimitterWithBuildIn(rateBucket)(GetSecInfoEndpoint)
	GetSecInfoEndpoint = kitzipkin.TraceEndpoint(localConfig.ZipKinTracer, "sec-info")(GetSecInfoEndpoint)

	//SecInfoListEndpoint

	GetSecInfoListEndpoint := endpoint.MakeSecInfoListEndpoint(skAppService)
	GetSecInfoListEndpoint = plugins.NewTokenBuketLimitterWithBuildIn(rateBucket)(GetSecInfoListEndpoint)
	GetSecInfoListEndpoint = kitzipkin.TraceEndpoint(localConfig.ZipKinTracer, "sec-info-list")(GetSecInfoListEndpoint)

	/**
	秒杀接口单独限流
	*/
	secRateBucket := rate.NewLimiter(rate.Every(time.Microsecond*100), 1000)

	SecKillEndpoint := endpoint.MakeSecKillEndpoint(skAppService)
	SecKillEndpoint = plugins.NewTokenBuketLimitterWithBuildIn(secRateBucket)(SecKillEndpoint)
	SecKillEndpoint = kitzipkin.TraceEndpoint(localConfig.ZipKinTracer, "sec-kill")(SecKillEndpoint)

	testEndpoint := endpoint.MakeTestEndpoint(skAppService)
	testEndpoint = kitzipkin.TraceEndpoint(localConfig.ZipKinTracer, "sec-test")(testEndpoint)

	endpts := endpoint.SkAppEndpoints{
		SecKillEndpoint:        SecKillEndpoint,
		HealthCheckEndpoint:    healthCheckEndpoint,
		GetSecInfoEndpoint:     GetSecInfoEndpoint,
		GetSecInfoListEndpoint: GetSecInfoListEndpoint,
		TestEndpoint:           testEndpoint,
	}

	ctx := context.Background()
	//创建http.handler
	r := transport.MakeHttpHandler(ctx, endpts, localConfig.ZipKinTracer, localConfig.Logger)
	//http server
	go func() {
		//注册服务
		register.Register()
		handler := r
		errChan <- http.ListenAndServe(":"+servicePort, handler)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	error := <-errChan
	//服务注销
	register.Deregister()
	fmt.Println(error)
}
