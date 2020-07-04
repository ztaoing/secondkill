/**
* @Author:zhoutao
* @Date:2020/7/4 下午5:32
 */

package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"net/http"
	"secondkill/oauth-service/model"
	"secondkill/oauth-service/service"
)

var (
	OAuth2DetailsKey       = "OAuth2Details"
	OAuth2ClientDetailsKey = "OAuth2ClientDetails"
	OAuth2ErrorKey         = "OAuthError"
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
