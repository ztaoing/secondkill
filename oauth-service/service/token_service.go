/**
* @Author:zhoutao
* @Date:2020/7/3 上午9:48
 */

package service

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"secondkill/oauth-service/model"
)

var (
	ErrNotSupportGrantType = errors.New("not supported grant type")
)

type TokenGranter interface {
	Grant(ctx context.Context, grantType string, client *model.ClientDetails, r *http.Request) (*model.OAuth2Token, error)
}

/**
组成token字典
注意：使用组合模式，使得不同的授权类型使用不同的tokenGrant接口实现结构体来生成访问令牌
组合节点ComposeTokenGranter
*/
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

//token增强
type TokenEnhancer interface {
	//组装Token信息
	Enhance(oauth2Token *model.OAuth2Token, oauth2Details *model.OAuth2Details) (*model.OAuth2Token, error)
	//从Token中还原信息
	Extract(tokenValue string) (*model.OAuth2Token, *model.OAuth2Details, error)
}

/**
令牌管理 + token增强
*/
type DefaultTokenService struct {
	tokenStore    TokenStore
	tokenEnhancer TokenEnhancer //token增强
}

func (d *DefaultTokenService) GetOAuth2DetailsByAccessToken(tokenValue string) (*model.OAuth2Details, error) {
	panic("implement me")
}

func (d *DefaultTokenService) CreateAccessToken(oauth2Details *model.OAuth2Details) (*model.OAuth2Token, error) {
	panic("implement me")
}

func (d *DefaultTokenService) RefreshAccessToken(refreshTokenValue string) (*model.OAuth2Token, error) {
	panic("implement me")
}

func (d *DefaultTokenService) GetAccessToken(details *model.OAuth2Token) (*model.OAuth2Token, error) {
	panic("implement me")
}

func (d *DefaultTokenService) ReadAccessToken(TokenValue string) (*model.OAuth2Token, error) {
	panic("implement me")
}

func NewTokenService(tokenStore TokenStore, enhancer TokenEnhancer) TokenService {
	return &DefaultTokenService{
		tokenStore:    tokenStore,
		tokenEnhancer: enhancer,
	}
}

/**
jsonWebToken 的token增强
*/
type JwtTokenEnhancer struct {
	secretKey []byte
}

func (enhance *JwtTokenEnhancer) Enhance(oauth2Token *model.OAuth2Token, oauth2Details *model.OAuth2Details) (*model.OAuth2Token, error) {
	panic("implement me")
}

func (enhance *JwtTokenEnhancer) Extract(tokenValue string) (*model.OAuth2Token, *model.OAuth2Details, error) {
	panic("implement me")
}

//将令牌对应的用户信息和客户端信息写入到JWT的声明中
func (enhance *JwtTokenEnhancer) sign(token *model.OAuth2Token, details *model.OAuth2Details) (*model.OAuth2Token, error) {

}

func NewJwtTokenEnhancer(secretKey string) TokenEnhancer {
	return &JwtTokenEnhancer{
		secretKey: []byte(secretKey),
	}
}

/**
jsonWebToken的令牌存储器
*/
type JwtTokenStore struct {
	jwtTokenEnhancer *JwtTokenEnhancer
}

func (j *JwtTokenStore) StoreAccessToken(oauth2Token *model.OAuth2Token, oauth2Details *model.OAuth2Details) {
	panic("implement me")
}

func (j *JwtTokenStore) ReadAccessToken(tokenValue string) (*model.OAuth2Token, error) {
	panic("implement me")
}

func (j *JwtTokenStore) ReadOAuth2Details(tokenValue string) (*model.OAuth2Details, error) {
	panic("implement me")
}

func (j *JwtTokenStore) GetAccessToken(oath2Details *model.OAuth2Details) (*model.OAuth2Token, error) {
	panic("implement me")
}

func (j *JwtTokenStore) RemoveAccessToken(tokenValue string) {
	panic("implement me")
}

func (j *JwtTokenStore) StoreRefreshToken(oauth2Token *model.OAuth2Token, oauth2Details *model.OAuth2Details) {
	panic("implement me")
}

func (j *JwtTokenStore) RemoveRefreshToken(oauth2Token string) {
	panic("implement me")
}

func (j *JwtTokenStore) ReadRefreshToken(tokenValue string) (*model.OAuth2Token, error) {
	panic("implement me")
}

func (j *JwtTokenStore) ReadOAuth2DetailsForRefreshToken(tokenValue string) (*model.OAuth2Details, error) {
	panic("implement me")
}

func NewJwtTokenStore(enhancer *JwtTokenEnhancer) TokenStore {
	return &JwtTokenStore{
		jwtTokenEnhancer: enhancer,
	}
}

//声明信息
type OAuth2TokenCustomClaims struct {
	//用户信息
	UserDetails model.UserDetails
	//客户端信息
	ClientDetails model.ClientDetails
	//用于刷新的令牌
	RefreshToken model.OAuth2Token
	//jwt标准声明
	jwt.StandardClaims
}
