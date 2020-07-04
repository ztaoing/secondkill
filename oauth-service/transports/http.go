/**
* @Author:zhoutao
* @Date:2020/7/4 下午8:28
* http方式的transport
 */

package transports

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/zipkin"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	zipkinGo "github.com/openzipkin/zipkin-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"secondkill/oauth-service/endpoint"
	"secondkill/oauth-service/service"
)

func MakeHttpHandler(ctx context.Context, endpoints endpoint.OAuthEndpoints, service service.TokenService, ClientdetailsService service.ClientDetailsService, tracer *zipkinGo.Tracer, logger log.Logger) http.Handler {
	r := mux.NewRouter()

	zipkinServer := zipkin.HTTPServerTrace(tracer, zipkin.Name("transport with http"))

	option := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
		zipkinServer,
	}
	//promhttp.Handler() return InstrumentMetricHandler
	r.Path("/metrics").Handler(promhttp.Handler())

	clientAuthorizationOptions := []kithttp.ServerOption{
		kithttp.ServerBefore(makeClientAuthorizationContext(ClientdetailsService, logger)),
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
		zipkinServer,
	}
	//kithttp.NewServer()，server实现了ServeHTTP
	r.Methods("POST").Path("/oauth/token").Handler(kithttp.NewServer(
		endpoints.TokenEndpoint,
		decodeCheckTokenRequest,
		encodeCheckTokenResponse,
		clientAuthorizationOptions...,
	))

	r.Methods("POST").Path("/oauth/check_token").Handler(kithttp.NewServer(
		endpoints.CheckTokenEndpoint,
		decodeCheckTokenRequest,
		encodeJsonResponse,
		clientAuthorizationOptions...,
	))

	//此请求之前不需要对客户端进行校验
	r.Methods("GET").Path("/health").Handler(kithttp.NewServer(
		endpoints.HealthCheckEndpoint,
		decodeHealthCheckRequest,
		encodeJsonResponse,
		option...,
	))
}

/**
ServerErrorEncoder用于在处理请求时遇到错误时对http.ResponseWriter进行编码。 客户可以使用它来提供自定义错误格式和响应代码
*/
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

//在之前处理：校验客户端权限
func makeClientAuthorizationContext(ClientDetailsService service.ClientDetailsService, logger log.Logger) kithttp.RequestFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		if clientName, clientPassword, ok := r.BasicAuth(); ok {
			clientDetails, err := ClientDetailsService.GetClientDetailsByClientId(ctx, clientName, clientPassword)
			if err == nil {
				return context.WithValue(ctx, endpoint.OAuth2ClientDetailsKey, clientDetails)
			}
		}
		return context.WithValue(ctx, endpoint.OAuth2ClientDetailsKey, endpoint.ErrInvalidClientRequest)
	}
}
