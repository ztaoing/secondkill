/**
* @Author:zhoutao
* @Date:2020/7/6 上午6:44
 */

package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"secondkill/sk-app/model"
	"secondkill/sk-app/service"
)

type SkAppEndpoints struct {
	SecKillEndpoint        endpoint.Endpoint
	HealthCheckEndpoint    endpoint.Endpoint
	GetSecInfoEndpoint     endpoint.Endpoint
	GetSecInfoListEndpoint endpoint.Endpoint
	TestEndpoint           endpoint.Endpoint
}

type SecInfoRequest struct {
	productId int `json:"product_id"`
}

type Response struct {
	Result map[string]interface{} `json:"result"`
	Error  error                  `json:"error"`
	Code   int                    `json:"code"`
}

type SecInfoListResponse struct {
	Result []map[string]interface{} `json:"result"`
	Error  error
	Code   int
}

/**
SecInfoEndpoint
*/
func MakeSecInfoEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(SecInfoRequest)

		resp := svc.SecInfo(req.productId)
		return Response{
			Result: resp,
			Error:  nil,
		}, nil
	}
}

/**
SecInfoListEndpoint
*/
func MakeSecInfoListEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		resp, num, err := svc.SecInfoList()
		return SecInfoListResponse{
			Result: resp,
			Error:  err,
			Code:   num,
		}, nil
	}
}

/**
SecKillEndpoint
*/
func MakeSecKillEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(model.SecRequest)
		result, code, err := svc.SecKill(&req)
		return Response{
			Result: result,
			Code:   code,
			Error:  err,
		}, nil
	}
}

/**
TestEndpoint
*/
func MakeTestEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return Response{
			Result: nil,
			Code:   1,
			Error:  nil,
		}, nil
	}
}

/**
健康检查的endpoint
*/
type HealthRequest struct {
}

type HealthResponse struct {
	Status bool `json:"status"`
}

func MakeHealthCheckEndpoint(service service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		status := service.HealthCheck()
		return HealthResponse{
			Status: status,
		}, nil
	}
}
