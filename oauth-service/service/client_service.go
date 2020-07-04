/**
* @Author:zhoutao
* @Date:2020/7/3 上午9:48
* 客户端信息查询
 */

package service

import (
	"context"
	"errors"
	"secondkill/oauth-service/model"
)

var ErrClientMessage = errors.New("invalid client")

type ClientDetailsService interface {
	GetClientDetailsByClientId(ctx context.Context, clientId string, clientSecret string) (*model.ClientDetails, error)
}

type MysqlClientDetailsService struct {
}

func NewMysqlClientDetailsService() ClientDetailsService {
	return &MysqlClientDetailsService{}
}

func (m *MysqlClientDetailsService) GetClientDetailsByClientId(ctx context.Context, clientId string, clientSecret string) (*model.ClientDetails, error) {
	clientDetailsModel := model.NewClientDetailsModel()

	if clientDetails, err := clientDetailsModel.GetClientDetailsByClient(clientId); err == nil {
		//正常获取
		if clientSecret == clientDetails.ClientSecret {
			return clientDetails, err
		} else {
			return nil, ErrClientMessage
		}

	} else {
		return nil, ErrClientMessage
	}
}
