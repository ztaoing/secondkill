/**
* @Author:zhoutao
* @Date:2020/7/1 上午7:16
 */

package client

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"secondkill/pb"
	"secondkill/pkg/discover"
	"secondkill/pkg/loadbalance"
)

type UserClient interface {
	CheckUser(ctx context.Context, tracer opentracing.Tracer, request *pb.UserRequest) (*pb.UserResponse, error)
}

/**
* 可以配置负载均衡策略，重试等机制。也可以配置invokeBefore和invokeAfter
 */
type UserClientImpl struct {
	manager     ClientManager
	serviceName string
	loadBalance loadbalance.LoadBalance
	tracer      opentracing.Tracer
}

func (impl *UserClientImpl) CheckUser(ctx context.Context, tracer opentracing.Tracer, request *pb.UserRequest) (*pb.UserResponse, error) {
	response := new(pb.UserResponse)
	///rpc调用地址：pb.UserService/Check
	if err := impl.manager.DecoratorInvoke("/pb.UserService/Check", "user_check", tracer, ctx, request, response); err == nil {
		return response, nil
	} else {
		return nil, err
	}

}

func NewUserClient(serviceName string, lb loadbalance.LoadBalance, tracer opentracing.Tracer) (UserClient, error) {
	if serviceName == "" {
		serviceName = "user"
	}
	if lb == nil {
		lb = defaultLoadBalance
	}

	return &UserClientImpl{
		manager: &DefaultClientManager{
			serviceName:     serviceName,
			loadbalance:     lb,
			discoveryClient: discover.ConsulService,
			logger:          discover.Logger,
		},
		serviceName: serviceName,
		loadBalance: lb,
		tracer:      tracer,
	}, nil
}
