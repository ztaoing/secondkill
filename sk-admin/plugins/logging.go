/**
* @Author:zhoutao
* @Date:2020/7/7 下午5:40
 */

package plugins

import (
	"github.com/go-kit/kit/log"
	"github.com/gohouse/gorose/v2"
	"secondkill/sk-admin/model"
	"secondkill/sk-admin/service"
	"time"
)

//基础service
type skAdminLoggingMiddleware struct {
	service.Service
	logger log.Logger
}

//activity
type activityLoggingMiddleware struct {
	service.ActivityService
	logger log.Logger
}

//product
type productLoggingMiddleware struct {
	service.ProductService
	logger log.Logger
}

func SKAdminLoggingMiddleware(logger log.Logger) service.ServiceMiddleware {
	return func(next service.Service) service.Service {
		return skAdminLoggingMiddleware{
			next, logger,
		}
	}
}

func ActivityLoggingMiddleware(logger log.Logger) service.ActivityServiceMiddleware {
	return func(next service.ActivityService) service.ActivityService {
		return activityLoggingMiddleware{
			next, logger,
		}
	}
}

func ProductLoggingMiddleware(logger log.Logger) service.ProductServiceMiddleware {
	return func(next service.ProductService) service.ProductService {
		return productLoggingMiddleware{
			next, logger,
		}
	}
}

func (s skAdminLoggingMiddleware) HealthCheck() (result bool) {
	defer func(begin time.Time) {
		s.logger.Log(
			"function", "HealthCheck",
			"result", result,
			"took", time.Since(begin),
		)
	}(time.Now())

	result = s.Service.HealthCheck()
	return
}

func (a activityLoggingMiddleware) GetActivityList() ([]gorose.Data, error) {
	defer func(begin time.Time) {
		a.logger.Log(
			"function", "GetActivityList",
			"took", time.Since(begin),
		)
	}(time.Now())

	ActivityList, err := a.ActivityService.GetActivityList()
	return ActivityList, err
}

func (a activityLoggingMiddleware) CreateActivity(activity *model.Activity) error {
	defer func(begin time.Time) {
		a.logger.Log(
			"function", "CreateActivity",
			"activity", activity,
			"took", time.Since(begin),
		)
	}(time.Now())

	err := a.ActivityService.CreateActivity(activity)
	return err
}

func (p productLoggingMiddleware) CreateProduct(product *model.Product) (err error) {
	defer func(begin time.Time) {
		p.logger.Log(
			"function", "CreateProduct",
			"product", product,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = p.ProductService.CreateProduct(product)
	return
}

func (p productLoggingMiddleware) GetProductList() ([]gorose.Data, error) {
	defer func(begin time.Time) {
		p.logger.Log(
			"function", "GetProductList",
			"took", time.Since(begin),
		)
	}(time.Now())

	data, err := p.ProductService.GetProductList()
	return data, err
}
