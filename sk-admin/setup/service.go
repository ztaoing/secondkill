/**
* @Author:zhoutao
* @Date:2020/7/7 下午9:34
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
	"net/http"
	"os"
	"os/signal"
	"secondkill/oauth-service/config"
	pkgConfig "secondkill/pkg/config"
	register "secondkill/pkg/discover"
	"secondkill/sk-admin/endpoint"
	"secondkill/sk-admin/plugins"
	"secondkill/sk-admin/service"
	"secondkill/sk-admin/transport"
	"syscall"
	"time"
)

//初始化HTTP服务
func InitHTTP(host, servicePort string) {
	flag.Parse()

	errChan := make(chan error)

	fieldKeys := []string{"method"}
	requestCount := kitPrometheus.NewCounterFrom(stdPrometheus.CounterOpts{
		Namespace: "seckill",
		Subsystem: "user_service",
		Name:      "request_count",
		Help:      "number of request received",
	}, fieldKeys)

	requestLatency := kitPrometheus.NewSummaryFrom(stdPrometheus.SummaryOpts{
		Namespace: "seckill",
		Subsystem: "user_service",
		Name:      "request_lantency",
		Help:      "Total duration of requests in microseconds",
	}, fieldKeys)

	//令牌桶
	ratebucket := rate.NewLimiter(rate.Every(time.Second*1), 100)

	var (
		skAdminService  service.Service
		activityService service.ActivityService
		productService  service.ProductService
	)

	skAdminService = service.SKAdminService{}
	activityService = service.ActivityServiceImpl{}
	productService = service.ProductServiceImpl{}

	//skAdminService
	skAdminService = plugins.SKAdminLoggingMiddleware(pkgConfig.Logger)(skAdminService)
	skAdminService = plugins.SkAdminMetrics(requestCount, requestLatency)(skAdminService)

	//activityService
	activityService = plugins.ActivityLoggingMiddleware(pkgConfig.Logger)(activityService)
	activityService = plugins.ActivityMetrics(requestCount, requestLatency)(activityService)

	//productService
	productService = plugins.ProductLoggingMiddleware(pkgConfig.Logger)(productService)
	productService = plugins.ProductMetrics(requestCount, requestLatency)(productService)

	//createActivityEnd
	createActivityEnd := endpoint.MakeCreateActivityEndpoint(activityService)
	//限流
	createActivityEnd = plugins.NewTokenBucketLimitterWithBuidIn(ratebucket)(createActivityEnd)
	//链路追踪
	createActivityEnd = kitzipkin.TraceEndpoint(config.ZipKinTracer, "create-activity")(createActivityEnd)

	//GetActivityEnd
	GetActivityEnd := endpoint.MakeGetActivityEndpoint(activityService)
	GetActivityEnd = plugins.NewTokenBucketLimitterWithBuidIn(ratebucket)(GetActivityEnd)
	GetActivityEnd = kitzipkin.TraceEndpoint(config.ZipKinTracer, "get-activity")(GetActivityEnd)

	//CreateProductEnd
	CreateProductEnd := endpoint.MakeCreateProductEndpoint(productService)
	CreateProductEnd = plugins.NewTokenBucketLimitterWithBuidIn(ratebucket)(CreateProductEnd)
	CreateProductEnd = kitzipkin.TraceEndpoint(config.ZipKinTracer, "create-product")(CreateProductEnd)

	//GetProductEnd
	GetProductEnd := endpoint.MakeGetProductEndpoint(productService)
	GetProductEnd = plugins.NewTokenBucketLimitterWithBuidIn(ratebucket)(GetProductEnd)
	GetProductEnd = kitzipkin.TraceEndpoint(config.ZipKinTracer, "get-product")(GetProductEnd)

	//healthCheck
	healthCheckEnd := endpoint.MakeHealthCheckEndpoint(skAdminService)
	healthCheckEnd = kitzipkin.TraceEndpoint(config.ZipKinTracer, "health-check")(healthCheckEnd)

	endpts := endpoint.SKAdminEndpoint{
		GetActivityEndpoint:    GetActivityEnd,
		CreateActivityEndpoint: createActivityEnd,

		GetProductEndpoint:    GetProductEnd,
		CreateProductEndpoint: CreateProductEnd,

		HealthCheckEndpoint: healthCheckEnd,
	}
	ctx := context.Background()

	//创建http.handler
	r := transport.MakeHttpHandler(ctx, endpts, config.ZipKinTracer, config.Logger)

	//http server
	go func() {
		register.Register()
		handler := r
		errChan <- http.ListenAndServe(":"+servicePort, handler)
	}()

	//停止
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	//注销服务
	error := <-errChan
	register.Deregister()
	fmt.Println(error)

}
