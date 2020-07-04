/**
* @Author:zhoutao
* @Date:2020/7/3 上午9:48
* 用户信息管理
 */

package service

import (
	"context"
	"errors"
	"secondkill/oauth-service/model"
	"secondkill/pb"
	"secondkill/pkg/client"
)

var (
	InvalidAuthentication = errors.New("invalid auth")
	InvalidUserInfo       = errors.New("invalid user info")
)

type UserDetailsService interface {
	GetUserDetailsByUsername(ctx context.Context, username string) (*model.UserDetails, error)
}

//实现了UserDetailsService接口
type RemoteUserService struct {
	userClient client.UserClient
}

//使用grpc 根据用户名获取用户信息
func (service *RemoteUserService) GetUserDetailsByUsername(ctx context.Context, username, password string) (*model.UserDetails, error) {
	resp, err := service.userClient.CheckUser(ctx, nil, &pb.UserRequest{
		Username: username,
		Password: password,
	})
	if err == nil {
		if resp.UserId != 0 {
			//成功获得用户信息
			return &model.UserDetails{
				UserId:   resp.UserId,
				UserName: username,
				Password: password,
			}, nil
		} else {
			return nil, InvalidAuthentication
		}
	} else {
		return nil, err
	}
}

func NewRemoteUserDetailsService() *RemoteUserService {
	userClient, _ := client.NewUserClient("user", nil, nil)
	return &RemoteUserService{
		userClient: userClient,
	}
}

//定义service中间件
type ServiceMiddleware func(Service) Service
