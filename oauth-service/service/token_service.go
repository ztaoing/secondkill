/**
* @Author:zhoutao
* @Date:2020/7/3 上午9:48
 */

package service

import (
	"context"
	"errors"
	"net/http"
	"secondkill/oauth-service/model"
)

var (
	ErrNotSupportGrantType = errors.New("not supported grant type")
)

type TokenGranter interface {
	Grant(ctx context.Context, grantType string, client *model.ClientDetails, r *http.Request) (*model.OAuth2Token, error)
}

type CommposeTokenGranter struct {
	//token授权:map[grantType]TokenGranter
	TokenGrantDict map[string]TokenGranter
}

//授权token
func (tokenGranter *CommposeTokenGranter) Grant(ctx context.Context, grantType string, client *model.ClientDetails, r *http.Request) (*model.OAuth2Token, error) {
	dispatchGranter := tokenGranter.TokenGrantDict[grantType]
	if dispatchGranter == nil {
		return nil, ErrNotSupportGrantType
	}
	return dispatchGranter.Grant(ctx, grantType, client, r)
}

func NewComposeTokenGranter(tokenGrantDict map[string]TokenGranter) TokenGranter {
	return &CommposeTokenGranter{
		TokenGrantDict: tokenGrantDict,
	}
}

/**
用户名密码Token授权
*/
type UsernamePasswordTokenGranter struct {
	supportGrantType   string
	userDetailsService UserDetailsService
	tokenService       TokenService //令牌管理
}

func (u UsernamePasswordTokenGranter) Grant(ctx context.Context, grantType string, client *model.ClientDetails, r *http.Request) (*model.OAuth2Token, error) {
	panic("implement me")
}

func NewUsernamePasswordTokenGranter(grantType string, userDetailsService UserDetailsService, tokenService TokenService) TokenGranter {
	return &UsernamePasswordTokenGranter{
		supportGrantType:   grantType,
		userDetailsService: userDetailsService,
		tokenService:       tokenService,
	}
}

/**
刷新令牌授权
*/
type RefreshTokenGranter struct {
	supportGrantType string
	tokenService     TokenService //令牌管理
}

func (r2 RefreshTokenGranter) Grant(ctx context.Context, grantType string, client *model.ClientDetails, r *http.Request) (*model.OAuth2Token, error) {
	panic("implement me")
}

func NewRefreshGrant(grantType string, tokenService TokenService) TokenGranter {
	return &RefreshTokenGranter{
		supportGrantType: grantType,
		tokenService:     tokenService,
	}
}

/**
令牌服务
用于令牌的管理
*/
type TokenService interface {
	//根据访问令牌获取对应的用户信息和客户端信息
	GetOAuth2DetailsByAccessToken(tokenValue string) (*model.OAuth2Details, error)
	//根据用户信息和客户端信息生成访问令牌
	CreateAccessToken(oauth2Details *model.OAuth2Details) (*model.OAuth2Token, error)
	//根据刷新令牌获取访问令牌
	RefreshAccessToken(refreshTokenValue string) (*model.OAuth2Token, error)
	//根据用户信息和客户端信息获取已生成访问令牌
	GetAccessToken(details *model.OAuth2Token) (*model.OAuth2Token, error)
	//根据访问令牌值获取访问令牌结构
	ReadAccessToken(TokenValue string) (*model.OAuth2Token, error)
}

/**
令牌存储器
负责存储生成的的令牌并维护令牌、用户、客户端之间的绑定关系
*/

type TokenStore interface {
	//存储访问令牌
	StoreAccessToken(oauth2Token *model.OAuth2Token, oauth2Details *model.OAuth2Details)
	//根据令牌值获取访问令牌结构
	ReadAccessToken(tokenValue string) (*model.OAuth2Token, error)
	//根据令牌值获取客户端和用户信息
	ReadOAuth2Details(tokenValue string) (*model.OAuth2Details, error)
	//根据客户端和用户信息获取访问令牌
	GetAccessToken(oath2Details *model.OAuth2Details) (*model.OAuth2Token, error)
	//移除访问令牌
	RemoveAccessToken(tokenValue string)
	//存储刷新令牌
	StoreRefreshToken(oauth2Token *model.OAuth2Token, oauth2Details *model.OAuth2Details)
	//移除刷新令牌
	RemoveRefreshToken(oauth2Token string)
	//根据令牌值获取刷新令牌
	ReadRefreshToken(tokenValue string) (*model.OAuth2Token, error)
	//根据令牌值获取令牌绑定的客户端和用户信息
	ReadOAuth2DetailsForRefreshToken(tokenValue string) (*model.OAuth2Details, error)
}
