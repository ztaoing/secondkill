/**
* @Author:zhoutao
* @Date:2020/7/7 下午8:46
 */

package transport

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/zipkin"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	gozipkin "github.com/openzipkin/zipkin-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"secondkill/sk-admin/endpoint"
	"secondkill/sk-admin/model"
	endpts "secondkill/sk-app/endpoint"
)

func MakeHttpHandler(ctx context.Context, endpoints endpoint.SKAdminEndpoint, zipkinTracer *gozipkin.Tracer, logger log.Logger) http.Handler {
	r := mux.NewRouter()

	zipkinServer := zipkin.HTTPServerTrace(zipkinTracer, zipkin.Name("http-transport"))

	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
		zipkinServer,
	}

	r.Methods("GET").Path("/product/list").Handler(kithttp.NewServer(
		endpoints.GetProductEndpoint,
		decodeGetListRequest,
		encodeResponse,
		options...,
	))

	r.Methods("POST").Path("/product/create").Handler(kithttp.NewServer(
		endpoints.CreateProductEndpoint,
		decodeCreateProductRequest,
		encodeResponse,
		options...,
	))
	r.Methods("POST").Path("/activity/create").Handler(kithttp.NewServer(
		endpoints.CreateActivityEndpoint,
		decodeCreateActivityRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/activity/list").Handler(kithttp.NewServer(
		endpoints.GetActivityEndpoint,
		decodeGetListRequest,
		encodeResponse,
		options...,
	))

	r.Path("/metrics").Handler(promhttp.Handler())

	r.Methods("GET").Path("/healht").Handler(kithttp.NewServer(
		endpoints.HealthCheckEndpoint,
		decodeHealthCheckRequest,
		encodeResponse,
		options...,
	))

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	return loggedRouter
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func decodeHealthCheckRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return endpts.HealthRequest{}, nil
}

func decodeCreateProductRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var product model.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		return nil, err
	}
	return product, nil
}
func decodeCreateActivityRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var activity model.Activity
	if err := json.NewDecoder(r.Body).Decode(&activity); err != nil {
		return nil, err
	}
	return activity, nil
}
