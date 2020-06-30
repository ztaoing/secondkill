/**
* @Author:zhoutao
* @Date:2020/6/30 下午10:48
 */

package client

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"secondkill/pb"
)

//我们以Auth服务提供的rpc接口为例，来介绍rpc客户端装饰器组件是如何运转的

type OAuthClient interface {
	CheckToken(ctx context.Context, tracer opentracing.Tracer, request *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error)
}
