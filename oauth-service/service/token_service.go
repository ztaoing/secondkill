/**
* @Author:zhoutao
* @Date:2020/7/3 上午9:48
 */

package service

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"secondkill/oauth-service/model"
	"strconv"
	"time"
)

var (
	ErrNotSupportGrantType = errors.New("not supported grant type")
	ErrExpiredToken        = errors.New("expired token")
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
	GetAccessToken(details *model.OAuth2Details) (*model.OAuth2Token, error)
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
令牌管理 +存储管理+ token增强
*/
type DefaultTokenService struct {
	tokenStore    TokenStore
	tokenEnhancer TokenEnhancer //token增强
}

//根据访问令牌获取对应的用户信息和客户端信息
func (d *DefaultTokenService) GetOAuth2DetailsByAccessToken(tokenValue string) (*model.OAuth2Details, error) {
	//首先判断是否过期
	accessToken, err := d.tokenStore.ReadAccessToken(tokenValue)
	if err == nil {
		if accessToken.IsExpired() {
			return nil, ErrExpiredToken
		}
		return d.tokenStore.ReadOAuth2Details(tokenValue)
	}
	return nil, err

}

/**
生成令牌：访问令牌和刷新令牌
*/
func (d *DefaultTokenService) CreateAccessToken(oauth2Details *model.OAuth2Details) (*model.OAuth2Token, error) {
	//通过客户端信息和用户信息查询令牌是否已经生成
	token, err := d.tokenStore.GetAccessToken(oauth2Details)
	var refreshToken *model.OAuth2Token
	if err == nil {
		//存在未失效的访问令牌，直接返回
		if !token.IsExpired() {
			d.tokenStore.StoreAccessToken(token, oauth2Details)
			return token, nil
		}
		//令牌已过期,移除访问令牌
		d.tokenStore.RemoveAccessToken(token.TokenValue)
		//刷新令牌存在时一起移除
		if token.RefreshToken != nil {
			d.tokenStore.RemoveRefreshToken(token.RefreshToken.TokenType)
		}
	}

	//令牌不存在

	//创建刷新令牌
	if refreshToken == nil || refreshToken.IsExpired() {
		refreshToken, err = d.createRefreshToken(oauth2Details)
		if err != nil {
			return nil, err
		}
	}

	//创建访问令牌
	accessToken, err := d.createAccessToken(refreshToken, oauth2Details)
	if err == nil {
		//保存 访问令牌和刷新令牌
		d.tokenStore.StoreAccessToken(accessToken, oauth2Details)
		d.tokenStore.StoreRefreshToken(accessToken, oauth2Details)
	}

	return accessToken, err
}

//生成刷新令牌
func (d *DefaultTokenService) createRefreshToken(details *model.OAuth2Details) (*model.OAuth2Token, error) {
	validitySecond := details.Client.RefreshTokenValiditySeconds
	s, _ := time.ParseDuration(strconv.Itoa(validitySecond) + "s")
	expiredTime := time.Now().Add(s)
	refreshToken := &model.OAuth2Token{
		ExpireTime: &expiredTime,
		TokenValue: uuid.NewV4().String(),
	}

	if d.tokenEnhancer != nil {
		return d.tokenEnhancer.Enhance(refreshToken, details)
	}
	return refreshToken, nil
}

//生成访问令牌
func (d *DefaultTokenService) createAccessToken(refreshToken *model.OAuth2Token, details *model.OAuth2Details) (*model.OAuth2Token, error) {
	validitySeconds := details.Client.AccessTokenValiditySeconds
	s, _ := time.ParseDuration(strconv.Itoa(validitySeconds) + "s")
	expiredTime := time.Now().Add(s)
	accessToken := &model.OAuth2Token{
		RefreshToken: refreshToken,
		ExpireTime:   &expiredTime,
		TokenValue:   uuid.NewV4().String(),
	}
	if d.tokenEnhancer != nil {
		return d.tokenEnhancer.Enhance(accessToken, details)
	}
	return accessToken, nil
}

//根据刷新令牌获取访问令牌
func (d *DefaultTokenService) RefreshAccessToken(refreshTokenValue string) (*model.OAuth2Token, error) {
	refreshToken, err := d.tokenStore.ReadRefreshToken(refreshTokenValue)
	if err == nil {
		if refreshToken.IsExpired() {
			return nil, ErrExpiredToken
		}
		oauthDetails, err := d.tokenStore.ReadOAuth2DetailsForRefreshToken(refreshTokenValue)
		if err == nil {
			token, err := d.tokenStore.GetAccessToken(oauthDetails)
			if err == nil {
				//移除原有的访问令牌
				d.tokenStore.RemoveAccessToken(token.TokenValue)
			}

			//移除已使用的刷新令牌
			d.tokenStore.RemoveRefreshToken(refreshTokenValue)

			//创建刷新令牌
			refreshToken, err := d.createRefreshToken(oauthDetails)
			if err == nil {
				accessToken, err := d.createAccessToken(refreshToken, oauthDetails)
				if err == nil {
					//存储
					d.tokenStore.StoreRefreshToken(refreshToken, oauthDetails)
					d.tokenStore.StoreAccessToken(accessToken, oauthDetails)
				}
				return accessToken, nil
			}

		}
	}
	return nil, err
}

//根据客户端和用户信息获取访问令牌
func (d *DefaultTokenService) GetAccessToken(details *model.OAuth2Details) (*model.OAuth2Token, error) {
	return d.tokenStore.GetAccessToken(details)
}

//根据访问令牌值获取访问令牌结构
func (d *DefaultTokenService) ReadAccessToken(TokenValue string) (*model.OAuth2Token, error) {
	return d.tokenStore.ReadAccessToken(TokenValue)
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
	//封装
	return enhance.sign(oauth2Token, oauth2Details)
}

//解封装
func (enhance *JwtTokenEnhancer) Extract(tokenValue string) (*model.OAuth2Token, *model.OAuth2Details, error) {
	token, err := jwt.ParseWithClaims(tokenValue, &OAuth2TokenCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return enhance.secretKey, nil
	})
	if err == nil {
		claims := token.Claims.(*OAuth2TokenCustomClaims)
		expiredTime := time.Unix(claims.ExpiresAt, 0)
		return &model.OAuth2Token{
				RefreshToken: &claims.RefreshToken,
				TokenValue:   tokenValue,
				ExpireTime:   &expiredTime,
			}, &model.OAuth2Details{
				User:   &claims.UserDetails,
				Client: &claims.ClientDetails,
			}, nil
	}
	return nil, nil, err
}

//将令牌对应的用户信息和客户端信息写入到JWT的声明中
func (enhance *JwtTokenEnhancer) sign(token *model.OAuth2Token, details *model.OAuth2Details) (*model.OAuth2Token, error) {
	expiredTime := token.ExpireTime
	clientDetails := *details.Client
	userDetails := *details.User
	clientDetails.ClientSecret = ""
	userDetails.Password = ""

	claims := OAuth2TokenCustomClaims{
		UserDetails:   userDetails,
		ClientDetails: clientDetails,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredTime.Unix(),
			Issuer:    "System",
		},
	}

	if token.RefreshToken != nil {
		claims.RefreshToken = *token.RefreshToken
	}
	//HS256使用对称加密算法
	EncodeToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenValue, err := EncodeToken.SignedString(enhance.secretKey)
	if err == nil {
		token.TokenValue = tokenValue
		token.TokenType = "jwt"
		return token, nil
	}
	return nil, err

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

func (tokenStore *JwtTokenStore) StoreAccessToken(oauth2Token *model.OAuth2Token, oauth2Details *model.OAuth2Details) {
	panic("implement me")
}

func (tokenStore *JwtTokenStore) ReadAccessToken(tokenValue string) (*model.OAuth2Token, error) {
	token, _, err := tokenStore.jwtTokenEnhancer.Extract(tokenValue)
	return token, err
}

func (tokenStore *JwtTokenStore) ReadOAuth2Details(tokenValue string) (*model.OAuth2Details, error) {
	_, details, err := tokenStore.jwtTokenEnhancer.Extract(tokenValue)
	return details, err
}

func (tokenStore *JwtTokenStore) GetAccessToken(oath2Details *model.OAuth2Details) (*model.OAuth2Token, error) {

}

func (tokenStore *JwtTokenStore) RemoveAccessToken(tokenValue string) {
	panic("implement me")
}

func (tokenStore *JwtTokenStore) StoreRefreshToken(oauth2Token *model.OAuth2Token, oauth2Details *model.OAuth2Details) {
	panic("implement me")
}

func (tokenStore *JwtTokenStore) RemoveRefreshToken(oauth2Token string) {
	panic("implement me")
}

func (tokenStore *JwtTokenStore) ReadRefreshToken(tokenValue string) (*model.OAuth2Token, error) {
	token, _, err := tokenStore.jwtTokenEnhancer.Extract(tokenValue)
	return token, err
}

//根据令牌值获取刷新令牌对应的客户端和用户信息
func (tokenStore *JwtTokenStore) ReadOAuth2DetailsForRefreshToken(tokenValue string) (*model.OAuth2Details, error) {
	_, details, err := tokenStore.jwtTokenEnhancer.Extract(tokenValue)
	return details, err
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
