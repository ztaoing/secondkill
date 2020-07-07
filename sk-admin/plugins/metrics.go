/**
* @Author:zhoutao
* @Date:2020/7/7 下午4:22
 */

package plugins

import (
	"github.com/go-kit/kit/metrics"
	"github.com/gohouse/gorose/v2"
	"secondkill/sk-admin/model"
	"secondkill/sk-admin/service"
	"time"
)

//metricMiddleware 定义监控中间件，嵌入service
type SKAdminMetricMiddelware struct {
	service.Service
	requestCount    metrics.Counter
	requestLantency metrics.Histogram
}

type activityMetricMiddle struct {
	service.ActivityService
	requestCount    metrics.Counter
	requestLantency metrics.Histogram
}

type productMetricMiddleware struct {
	service.ProductService
	requestCount    metrics.Counter
	requestLantency metrics.Histogram
}

//metrics封装监控方法
func SkAdminMetrics(requestCount metrics.Counter, requestLantency metrics.Histogram) service.ServiceMiddleware {
	return func(next service.Service) service.Service {
		return SKAdminMetricMiddelware{
			next,
			requestCount,
			requestLantency,
		}
	}
}

func ProductMetrics(requestCount metrics.Counter, requestLantency metrics.Histogram) service.ProductServiceMiddleware {
	return func(next service.ProductService) service.ProductService {
		return productMetricMiddleware{
			next,
			requestCount,
			requestLantency,
		}
	}
}

func ActivityMetrics(requestCount metrics.Counter, requestLantency metrics.Histogram) service.ActivityServiceMiddleware {
	return func(next service.ActivityService) service.ActivityService {
		return activityMetricMiddle{
			next,
			requestCount,
			requestLantency,
		}
	}
}

//service指标参数
func (sk SKAdminMetricMiddelware) HealthCheck(result bool) {
	defer func(begin time.Time) {
		lvs := []string{"method", "HealthCheck"}
		sk.requestCount.With(lvs...).Add(1)
		sk.requestLantency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	result = sk.Service.HealthCheck()
	return
}

//productMetric
func (p productMetricMiddleware) CreateProduct(product *model.Product) error {
	defer func(begin time.Time) {
		lvs := []string{"method", "CreateProduct"}
		p.requestCount.With(lvs...).Add(1)
		p.requestLantency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	error := p.ProductService.CreateProduct(product)
	return error
}

func (p productMetricMiddleware) GetProductList() ([]gorose.Data, error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetProductList"}
		p.requestCount.With(lvs...).Add(1)
		p.requestLantency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	ProductList, err := p.ProductService.GetProductList()
	return ProductList, err
}

func (a activityMetricMiddle) CreateActivity(activity *model.Activity) error {
	defer func(begin time.Time) {
		lvs := []string{"method", "CreateActivity"}
		a.requestCount.With(lvs...).Add(1)
		a.requestLantency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err := a.ActivityService.CreateActivity(activity)
	return err
}

func (a activityMetricMiddle) GetActivityList() ([]gorose.Data, error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetActivityList"}
		a.requestCount.With(lvs...).Add(1)
		a.requestLantency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	ActivityList, err := a.ActivityService.GetActivityList()
	return ActivityList, err
}
