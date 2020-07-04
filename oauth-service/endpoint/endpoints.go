/**
* @Author:zhoutao
* @Date:2020/7/4 下午5:32
 */

package endpoint

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"net/http"
	"secondkill/oauth-service/model"
	"secondkill/oauth-service/service"
)

const (
	OAuth2DetailsKey       = "OAuth2Details"
	OAuth2ClientDetailsKey = "OAuth2ClientDetails"
	OAuth2ErrorKey         = "OAuthError"
)

var (
	ErrInvalidClientRequest = errors.New("invalid client request")
	ErrNotPermit            = errors.New("not permit")
	ErrInvalidUserRequest   = errors.New("invalid user request")
)

//统一认证与授权的endpoint层
type OAuthEndpoints struct {
	//token
	TokenEndpoint endpoint.Endpoint
	//验证token
	CheckTokenEndpoint endpoint.Endpoint
	//GRPC验证token
	GRPCCheckTokenEndpoint endpoint.Endpoint
	//健康检查
	HealthCheckEndpoint endpoint.Endpoint
}

//令牌认证
//客户端验证中间件
//验证请求上下文中是否携带了客户端信息，如果请求中没有携带验证过的客户端信息，将直接返回错误给请求方
//在进入Endpoint之前统一验证context中的OAuthDetails是否存在
func MakeClientAuthorizationMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			//请求上下文是否存在错误
			if err, ok := ctx.Value(OAuth2ErrorKey).(error); ok {
				return nil, err
			}
			//验证客户端信息和用户信息是否存在，不存在则拒绝访问
			if _, ok := ctx.Value(OAuth2ClientDetailsKey).(*model.ClientDetails); !ok {
				return nil, ErrInvalidClientRequest
			}
			return next(ctx, request)
		}
	}
}

//访问资源服务器受保护的资源的端点，不仅需要请求中携带有效的访问令牌，
//还需要访问令牌对应的用户和客户端具备足够的权限
//在transport层中makeOAuth2AuthroizationContext请求处理器中获得了用户信息和客户端信息，
//可以根据他们具备的权限等级，判断是否具备访问点的权限
func MakeAuthorityAuthorizationMiddleware(authority string, logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if err, ok := ctx.Value(OAuth2ErrorKey).(error); ok {
				return nil, err
			}
			if details, ok := ctx.Value(OAuth2ClientDetailsKey).(*model.OAuth2Details); !ok {
				return nil, ErrInvalidClientRequest
			} else {
				for _, value := range details.User.Authorities {
					//权限检查
					if value == authority {
						return next(ctx, request)
					}
				}
			}
			return nil, ErrNotPermit

		}
	}
}

//在进入endpoint之前统一验证context中的OAuth2Details是否存在
func MakeOAuth2AuthorizationMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if err, ok := ctx.Value(OAuth2ErrorKey).(error); !ok {
				return nil, err
			}
			if _, ok := ctx.Value(OAuth2ClientDetailsKey).(*model.OAuth2Details); !ok {
				return nil, ErrInvalidUserRequest
			}
			return next(ctx, request)
		}
	}
}

/**
token request-response
*/
type TokenRequest struct {
	GrantType string
	Reader    *http.Request
}

type TokenResponse struct {
	AccessToken *model.OAuth2Token `json:"access_token"`
	Error       string             `json:"error"`
}

func MakeTokenEndpoint(svc service.TokenGranter, clientService service.ClientDetailsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*TokenRequest)
		token, err := svc.Grant(ctx, req.GrantType, ctx.Value(OAuth2ClientDetailsKey).(*model.ClientDetails), req.Reader)

		var errString = ""
		if err != nil {
			errString = err.Error()
		}
		return TokenResponse{AccessToken: token, Error: errString}, nil
	}
}

/**
checkToken request-response
*/

type CheckTokenRequest struct {
	Token         string
	ClientDetails *model.OAuth2Details
}

type CheckTokenResponse struct {
	OAuthDetails *model.OAuth2Details `json:"o_auth_details"`
	Error        string               `json:"error"`
}

func MakeCheckTokenEndpoint(svc service.TokenService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CheckTokenRequest)
		tokenDetails, err := svc.GetOAuth2DetailsByAccessToken(req.Token)

		var errString = ""
		if err != nil {
			errString = err.Error()
		}
		return CheckTokenResponse{
			OAuthDetails: tokenDetails,
			Error:        errString,
		}, nil

	}
}

/**
simple request-response
*/

type SimpleRequest struct {
}

type SimpleResponse struct {
	Result string `json:"result"`
	Error  string `json:"error"`
}

/**
health request-response
*/

type HealthRequest struct {
}
type HealthResponse struct {
	Status bool `json:"status"`
}

//创建健康检查的endpoint
func MakeHealthCheckEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		status := svc.HealthCheck()
		return HealthResponse{
			Status: status,
		}, nil
	}
}
