/**
* @Author:zhoutao
* @Date:2020/7/2 下午3:50
 */

package route

import (
	"context"
	"errors"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-kit/kit/log"
	"github.com/openzipkin/zipkin-go"
	zipkinhttpSvr "github.com/openzipkin/zipkin-go/middleware/http"
	"net/http"
	"net/http/httputil"
	"secondkill/gateway/config"
	"secondkill/pb"
	"secondkill/pkg/client"
	"secondkill/pkg/discover"
	"secondkill/pkg/loadbalance"
	"strings"
	"sync"
)

type HystrixRouter struct {
	svcMap      *sync.Map               //服务护理，存储已经通过hystrix监控的服务列表
	logger      log.Logger              //日志工具
	fallbackMsg string                  //回调消息
	tracer      *zipkin.Tracer          //服务追踪对象
	loadbalance loadbalance.LoadBalance //负载均衡
}

func Routes(zipkinTracer *zipkin.Tracer, fallbackMsg string, logger log.Logger) http.Handler {
	return &HystrixRouter{
		svcMap:      &sync.Map{},
		logger:      logger,
		fallbackMsg: fallbackMsg,
		tracer:      zipkinTracer,
		loadbalance: &loadbalance.RandomLoadBalancce{},
	}
}

func (router *HystrixRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//查询原始请求路径 如：/string-service/calculate/1/5
	reqPath := r.URL.Path
	router.logger.Log("reqPath", reqPath)

	//健康检查直接返回
	if reqPath == "/health" {
		w.WriteHeader(200)
		return
	}

	//检验是否为非法请求
	var err error
	if reqPath == "" || !preFilter(r) {
		err = errors.New("illegal request")
		w.WriteHeader(403)
		w.Write([]byte(err.Error()))
		return
	}

	//获取服务名称
	pathArray := strings.Split(reqPath, "/")
	serviceName := pathArray[1]

	//检查是否已经加入监控
	if _, ok := router.svcMap.Load(serviceName); !ok {
		//没有在键控中
		//把serviceName作为命名对象
		hystrix.ConfigureCommand(serviceName, hystrix.CommandConfig{
			Timeout: 1000,
		})
		router.svcMap.Store(serviceName, serviceName)
	}
	//hystrix开始
	err = hystrix.Do(serviceName, func() error {
		//服务发现并根据负载均衡获取一个实例
		serviceInstance, err := discover.DiscoveryService(serviceName)
		if err != nil {
			return err
		}

		director := func(req *http.Request) {
			//重新组织请求路径，
			destPath := strings.Join(pathArray[2:], "/")

			router.logger.Log("serviceId", serviceInstance.Host, serviceInstance.Port)

			//设置代理服务地址
			req.URL.Scheme = "http"
			req.URL.Host = fmt.Sprintf("%s:%d", serviceInstance.Host, serviceInstance.Port)
			req.URL.Path = "/" + destPath

		}

		var proxyError error = nil
		//为反向代理增加追踪逻辑，使用roudtrip代替默认的transport
		/**
		RoundTrip执行一个HTTP事务，返回所提供请求的响应。 RoundTrip不应尝试解释响应。
		特别是，如果RoundTrip获得响应，则无论响应的HTTP状态代码如何，它都必须返回err == nil。
		 如果失败，则应保留非null错误。同样，RoundTrip不应尝试处理更高级别的协议详细信息，例如重定向，
		身份验证或cookie.RoundTrip不应修改请求，除非消耗并关闭请求的正文。
		RoundTrip可能在单独的goroutine中读取请求的字段。来电者除非响应的响应，否则不应更改或重用请求
		身体已经关闭。RoundTrip必须始终关闭身体，包括错误，但根据具体实施情况，可能需要单独进行
		即使在RoundTrip返回后也可以使用goroutine。这意味着希望重用主体以用于后续请求的调用者必须安排在等待Close调用之后再这样做。
		请求的URL和标头字段必须初始化。
		*/
		roundTrip, _ := zipkinhttpSvr.NewTransport(router.tracer, zipkinhttpSvr.TransportTrace(true))

		//反向代理失败错误处理
		errorHandler := func(ew http.ResponseWriter, er *http.Request, err error) {
			proxyError = err
		}
		/**
		Director:必须具有将请求修改为要使用传输发送的新请求的功能。 然后将其响应复制回未经修改的原始客户端。Director返回后不得访问提供的请求。
		Transport:用来执行代理请求
		*/
		proxy := &httputil.ReverseProxy{
			Director:     director,
			Transport:    roundTrip,
			ErrorHandler: errorHandler,
		}
		proxy.ServeHTTP(w, r)

		return proxyError

	}, func(err error) error {
		//run执行失败，返回fallback信息
		router.logger.Log("fallback error description", err.Error())
		return errors.New(router.fallbackMsg)
	})
	//do执行失败，返回fallback信息
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}
}

//提前过滤
func preFilter(r *http.Request) bool {
	reqPath := r.URL.Path
	if reqPath == "" {
		return false
	}

	res := config.Match(reqPath)
	if res {
		return true
	}

	authToken := r.Header.Get("Authorization")
	if authToken == "" {
		return false
	}

	OAclient, _ := client.NewOAuthClient("oauth", nil, nil)
	resp, remoteErr := OAclient.CheckToken(context.Background(), nil, &pb.CheckTokenRequest{Token: authToken})

	if remoteErr != nil || resp == nil {
		return false
	} else {
		return true
	}
}
