/**
* @Author:zhoutao
* @Date:2020/7/4 下午8:29
 */

package transports

import (
	"context"
	endpoint2 "secondkill/oauth-service/endpoint"
	"secondkill/oauth-service/model"
	"secondkill/pb"
)

//向rpc-server发送前的加密操作
func EncodeGRPCCheckRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*endpoint2.CheckTokenRequest)
	return &pb.CheckTokenRequest{
		Token: req.Token,
	}, nil
}

//从rpc-client接收请求后的解密操作
func DecodeGRPCCheckRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.CheckTokenRequest)
	return &endpoint2.CheckTokenRequest{
		Token: req.Token,
	}, nil
}

//从rpc-server返回的结果的加密操作
func EncodeGRPCCheckResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(endpoint2.CheckTokenResponse)
	if resp.Error != "" {
		return &pb.CheckTokenResponse{
			IsValidToken: false,
			Err:          resp.Error,
		}, nil
	} else {
		return &pb.CheckTokenResponse{
			UserDetails: &pb.UserDetails{
				UserId:      resp.OAuthDetails.User.UserId,
				UserName:    resp.OAuthDetails.User.UserName,
				Authorities: resp.OAuthDetails.User.Authorities,
			},
			ClientDetails: &pb.ClientDetails{
				ClientId:                    resp.OAuthDetails.Client.ClientId,
				AccessTokenValiditySeconds:  int32(resp.OAuthDetails.Client.AccessTokenValiditySeconds),
				RefreshTokenValiditySeconds: int32(resp.OAuthDetails.Client.RefreshTokenValiditySeconds),
				AuthorizedGrantTypes:        resp.OAuthDetails.Client.AutorizedGrantTypes,
			},
		}, nil
	}

}

//从rpc-server接收应答后的解密操作
func DecodeGRPCCheckResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pb.CheckTokenResponse)
	if resp.Err != "" {
		return endpoint2.CheckTokenResponse{
			OAuthDetails: nil,
			Error:        resp.Err,
		}, nil
	} else {
		return &endpoint2.CheckTokenResponse{
			OAuthDetails: &model.OAuth2Details{
				User: &model.UserDetails{
					UserId:      resp.UserDetails.UserId,
					UserName:    resp.UserDetails.UserName,
					Authorities: resp.UserDetails.Authorities,
				},
				Client: &model.ClientDetails{
					ClientId:                    resp.ClientDetails.ClientId,
					AccessTokenValiditySeconds:  int(resp.ClientDetails.AccessTokenValiditySeconds),
					RefreshTokenValiditySeconds: int(resp.ClientDetails.RefreshTokenValiditySeconds),
					AutorizedGrantTypes:         resp.ClientDetails.AuthorizedGrantTypes,
				},
			},
		}, nil
	}

}
