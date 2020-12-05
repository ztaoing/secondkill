/**
* @Author:zhoutao
* @Date:2020/6/30 下午10:48
 */

package client

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"secondkill/pb"
	"secondkill/pkg/discover"
	"secondkill/pkg/loadbalance"
)

type OAuthClient interface {
	CheckToken(ctx context.Context, tracer opentracing.Tracer, request *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error)
}

type OAuthClientImpl struct {
	manager     ClientManager           //客户端管理器
	serviceName string                  //服务的名称
	loadbalance loadbalance.LoadBalance //负载均衡策略
	tracer      opentracing.Tracer      //链路追踪系统
}

//对用户token进行校验
func (O *OAuthClientImpl) CheckToken(ctx context.Context, tracer opentracing.Tracer, request *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error) {
	response := new(pb.CheckTokenResponse)
	if err := O.manager.DecoratorInvoke("pb.OAuthService/CheckToken", "token_check", tracer, ctx, request, response); err == nil {
		return response, nil
	} else {
		return nil, err
	}
}

//初始化OAuthClientImpl实例，并配置其各种属性
func NewOAuthClient(serviceName string, lb loadbalance.LoadBalance, tracer opentracing.Tracer) (OAuthClient, error) {
	if serviceName == "" {
		serviceName = "oauth"
	}
	if lb == nil {
		lb = defaultLoadBalance
	}
	return &OAuthClientImpl{
		manager: &DefaultClientManager{
			serviceName:     serviceName,
			loadbalance:     lb,
			discoveryClient: discover.ConsulService,
			logger:          discover.Logger,
		},
		serviceName: serviceName,
		loadbalance: lb,
		tracer:      tracer,
	}, nil

}
