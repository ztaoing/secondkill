/**
* @Author:zhoutao
* @Date:2020/7/1 上午7:00
 */

package client

import (
	"context"
	"errors"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	zipkinbridge "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"google.golang.org/grpc"
	"log"
	"secondkill/pkg/bootstrap"
	"secondkill/pkg/config"
	"secondkill/pkg/discover"
	"secondkill/pkg/loadbalance"
	"strconv"
)

//默认负载均衡策略是 随机负载均衡
var defaultLoadBalance loadbalance.LoadBalance = &loadbalance.RandomLoadBalancce{}
var ErrRPCService = errors.New("no rpc service")

//客户端管理器
type ClientManager interface {
	DecoratorInvoke(path string, hystrixName string, tracer opentracing.Tracer, ctx context.Context, inputVal interface{}, outVal interface{}) (err error)
}

type DefaultClientManager struct {
	serviceName     string
	logger          *log.Logger
	discoveryClient discover.DiscoveryClient
	loadbalance     loadbalance.LoadBalance
	before          []InvokeBeforeFunc //前 func list
	after           []InvokeAfterFunc  //后 func list

}

//之前调用
type InvokeAfterFunc func() (err error)

//之后调用
type InvokeBeforeFunc func() (err error)

/**
装饰器调用
@param path rpc请求路径
@param hystrixName 方法名称
@param tracer 链路追踪系统
@param ctx 上下文
@param inputVal 请求
@param outVal 响应
*/
func (manager *DefaultClientManager) DecoratorInvoke(path string, hystrixName string, tracer opentracing.Tracer, ctx context.Context, inputVal interface{}, outVal interface{}) (err error) {

	//进行发送rpc请求 前的统一回调处理
	for _, fn := range manager.before {
		if err = fn(); err != nil {
			return err
		}
	}

	//使用hystrix的Do方法构造对应的断路器保护
	if err = hystrix.Do(hystrixName, func() error {
		//服务发现
		instanceList := manager.discoveryClient.DiscoverServices(manager.serviceName, manager.logger)
		//负载均衡
		if instance, err := manager.loadbalance.SelectService(instanceList); err == nil {
			if instance.GrpcPort > 0 {

				//获取服务的rpc端口并选取的实例发送rpc请求
				if conn, err := grpc.DialContext(ctx, instance.Host+":"+strconv.Itoa(instance.GrpcPort), grpc.WithInsecure(),
					grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(genTracer(tracer),
						otgrpc.LogPayloads()))); err == nil {
					if err = conn.Invoke(ctx, path, inputVal, outVal); err != nil {
						return err
					}
				} else {
					return err
				}

			}
		} else {
			return ErrRPCService
		}
		return nil
	}, func(err error) error {
		return err
	}); err != nil {
		return err
	} else {
		//调用clientManager 的after回调函数
		for _, fn := range manager.after {
			if err := fn(); err != nil {
				return err
			}
		}
		return nil
	}
}

//生成tracer
func genTracer(tracer opentracing.Tracer) opentracing.Tracer {
	if tracer != nil {
		return tracer
	}
	//如果没有tracer,生成默认tracer
	zipkinUrl := "http://" + config.TraceConfig.Host + ":" + config.TraceConfig.Port + config.TraceConfig.Url

	zipkinRecodeUrl := bootstrap.HttpConfig.Host + ":" + bootstrap.HttpConfig.Port

	//new Collector
	collector, err := zipkinbridge.NewHTTPCollector(zipkinUrl)
	if err != nil {
		log.Fatalf("zipkin.NewHTTPCollector err:%v", err)
	}

	//new Recorder
	recorder := zipkinbridge.NewRecorder(collector, false, zipkinRecodeUrl, bootstrap.DiscoverConfig.ServiceName)

	//new Tracer
	resTracer, err := zipkinbridge.NewTracer(recorder, zipkinbridge.ClientServerSameSpan(true))
	if err != nil {
		log.Fatalf("zipkinbridge.NewTracer err:%v", err)
	}
	return resTracer

}
