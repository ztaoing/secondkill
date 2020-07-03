/**
* @Author:zhoutao
* @Date:2020/7/3 上午8:03
 */

package model

import (
	"encoding/json"
	"log"
	"secondkill/pkg/mysql"
)

//客户端信息
type ClientDetails struct {
	//client id
	ClientId string
	//client 秘钥
	ClientSecret string
	//访问令牌的有效时间，秒
	AccessTokenValiditySeconds int
	//刷新令牌的有效时间，秒
	RefreshTokenValiditySeconds int
	//重定向地址，授权码类型中使用
	RegisteredRefirectUri string
	//可以使用的授权码类型
	AutorizedGrantTypes []string
}

//校验秘钥
func (clientDetails *ClientDetails) IsMatch(clientId string, clientSecret string) bool {
	return clientId == clientDetails.ClientId && clientSecret == clientDetails.ClientSecret
}

type ClientDetailsModel struct {
}

func NewClientDetailsModel() *ClientDetailsModel {
	return &ClientDetailsModel{}
}

func (p *ClientDetailsModel) getTabName() string {
	return "client_details"
}

//根据client id获取client信息
func (p *ClientDetailsModel) GetClientDetailsByClient(clientId string) (*ClientDetails, error) {
	conn := mysql.DB()
	if result, err := conn.Table(p.getTabName()).Where(map[string]interface{}{"client_id": clientId}).First(); err == nil {

		var authorizedGrantTypes []string
		_ = json.Unmarshal([]byte(result["authorized_grant_types"].(string)), &authorizedGrantTypes)

		return &ClientDetails{
			ClientId:                    result["client_id"].(string),
			ClientSecret:                result["client_secret"].(string),
			AccessTokenValiditySeconds:  int(result["access_token_validity_seconds"].(int64)),
			RefreshTokenValiditySeconds: int(result["fresh_token_validity_seconds"].(int64)),
			RegisteredRefirectUri:       result["registered_redirect_url"].(string),
			AutorizedGrantTypes:         authorizedGrantTypes,
		}, nil
	} else {
		return nil, err
	}
}

//创建client信息
func (p *ClientDetailsModel) CreateClientDetails(details *ClientDetails) error {
	conn := mysql.DB()
	grantTypeString, _ := json.Marshal(details.AutorizedGrantTypes)

	_, err := conn.Table(p.getTabName()).Data(map[string]interface{}{
		"client_id":                     details.ClientId,
		"client_secret":                 details.ClientSecret,
		"access_token_validity_seconds": details.AccessTokenValiditySeconds,
		"refesh_token_validity_seconds": details.RegisteredRefirectUri,
		"registered_redirect_url":       details.RegisteredRefirectUri,
		"authorized_grant_types":        grantTypeString,
	}).Insert()
	if err != nil {
		log.Printf("Error：%v", err)
		return err
	}
	//创建成功
	return nil
}
