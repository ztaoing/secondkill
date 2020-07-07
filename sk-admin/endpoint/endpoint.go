/**
* @Author:zhoutao
* @Date:2020/7/7 下午6:41
 */

package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/gohouse/gorose/v2"
	"secondkill/sk-admin/model"
	"secondkill/sk-admin/service"
)

type SKAdminEndpoint struct {
	GetActivityEndpoint    endpoint.Endpoint
	CreateActivityEndpoint endpoint.Endpoint

	CreateProductEndpoint endpoint.Endpoint
	GetProductEndpoint    endpoint.Endpoint

	HealthCheckEndpoint endpoint.Endpoint
}

//user
type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResponse struct {
}

//response
type GetResponse struct {
	Result []gorose.Data `json:"result"`
	Error  error         `json:"error"`
}

type CreateResponse struct {
	Error error `json:"error"`
}

//GetActivityEndpoint
func MakeGetActivityEndpoint(svc service.ActivityService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		activityList, err := svc.GetActivityList()
		if err != nil {
			return GetResponse{
				Result: nil,
				Error:  err,
			}, nil
		}
		return GetResponse{
			Result: activityList,
			Error:  err,
		}, nil
	}
}

//CreateActivityEndpoint
func MakeCreateActivityEndpoint(svc service.ActivityService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(model.Activity)
		err = svc.CreateActivity(&req)
		return CreateResponse{
			Error: err,
		}, nil
	}
}

//GetProductEndpoint
func MakeGetProductEndpoint(svc service.ProductService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		data, err := svc.GetProductList()
		if err != nil {
			return GetResponse{
				Result: nil,
				Error:  err,
			}, nil
		}
		return GetResponse{
			Result: data,
			Error:  nil,
		}, nil
	}
}

func MakeCreateProductEndpoint(svc service.ProductService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(model.Product)
		err = svc.CreateProduct(&req)
		return CreateResponse{
			Error: err,
		}, nil
	}
}

type HealthCheckRequest struct {
}
type HealthCheckResponse struct {
	Status bool `json:"status"`
}

func MakeHealthCheckEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		status := svc.HealthCheck()
		return HealthCheckResponse{
			Status: status,
		}, nil
	}
}
