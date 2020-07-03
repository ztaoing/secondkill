/**
* @Author:zhoutao
* @Date:2020/7/3 上午7:58
令牌管理
*/

package model

import "time"

type OAuth2Token struct {
	//刷新令牌
	RefreshToken *OAuth2Token
	//令牌类型
	TokenType string
	//令牌
	TokenValue string
	//过期时间
	ExpireTime *time.Time
}

type OAuth2Details struct {
	Client *ClientDetails
	User   *UserDetails
}

//过期
func (oauth2Token *OAuth2Token) IsExpired() bool {
	return oauth2Token.ExpireTime != nil && oauth2Token.ExpireTime.Before(time.Now())
}
