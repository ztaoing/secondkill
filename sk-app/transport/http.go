/**
* @Author:zhoutao
* @Date:2020/7/6 上午7:18
 */

package transport

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/zipkin"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	gozipkin "github.com/openzipkin/zipkin-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	endpts "secondkill/sk-app/endpoint"
	"secondkill/sk-app/model"
)

func MakeHttpHandler(ctx context.Context, endpoint endpts.SkAppEndpoints, tracer *gozipkin.Tracer, logger log.Logger) http.Handler {
	r := mux.NewRouter()

	zipkinServerOptions := zipkin.HTTPServerTrace(tracer, zipkin.Name("http-transport"))
	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
		zipkinServerOptions,
	}

	r.Methods("GET").Path("/sec/info").Handler(kithttp.NewServer(
		endpoint.GetSecInfoEndpoint,
		decodeSecInfoRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/sec/list").Handler(kithttp.NewServer(
		endpoint.GetSecInfoListEndpoint,
		decodeSecInfoListRequest,
		encodeResponse,
		options...,
	))

	r.Methods("POST").Path("/sec/kill").Handler(kithttp.NewServer(
		endpoint.SecKillEndpoint,
		decodeSecKillRequest,
		encodeResponse,
		options...,
	))

	r.Methods("POST").Path("/sec/test").Handler(kithttp.NewServer(
		endpoint.TestEndpoint,
		decodeTestRequest,
		encodeResponse,
		options...,
	))

	r.Path("/metrics").Handler(promhttp.Handler())

	r.Methods("GET").Path("/health").Handler(kithttp.NewServer(
		endpoint.HealthCheckEndpoint,
		decodeHealthCheckRequest,
		encodeResponse,
		options...,
	))
	return r
}

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

func decodeSecInfoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var secInfoRequest endpts.SecInfoRequest
	if err := json.NewDecoder(r.Body).Decode(&secInfoRequest); err != nil {
		return nil, err
	}
	return secInfoRequest, nil

}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// decodeHealthCheckRequest decode request
func decodeHealthCheckRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return endpts.HealthRequest{}, nil
}

func decodeTestRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return endpts.HealthRequest{}, nil
}

func decodeSecInfoListRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeSecKillRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var secRequest model.SecRequest
	if err := json.NewDecoder(r.Body).Decode(&secRequest); err != nil {
		return nil, err
	}
	return secRequest, nil
}
