/**
* @Author:zhoutao
* @Date:2020/7/4 下午8:28
* 使用grpc方式 传输
 */

package transports

import (
	"context"
	"github.com/go-kit/kit/transport/grpc"
	endpoint2 "secondkill/oauth-service/endpoint"
	"secondkill/pb"
)

type grpcServer struct {
	checkTokenServer grpc.Handler
}

func (g *grpcServer) CheckToken(ctx context.Context, request *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error) {
	_, resp, err := g.checkTokenServer.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.CheckTokenResponse), nil
}

func NewGRPCServer(_ context.Context, endpoints endpoint2.OAuthEndpoints, option grpc.ServerOption) pb.OAuthServiceServer {
	return &grpcServer{
		checkTokenServer: grpc.NewServer(
			endpoints.GRPCCheckTokenEndpoint,
			DecodeGRPCCheckRequest,
			EncodeGRPCCheckResponse,
			option,
		),
	}
}
