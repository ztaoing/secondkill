/**
* @Author:zhoutao
* @Date:2020/7/6 上午11:14
 */

package plugins

import (
	"github.com/go-kit/kit/metrics"
	"secondkill/sk-app/model"
	"secondkill/sk-app/service"
)

//metrics中间件
type skAppMetricMiddleware struct {
	service        service.Service
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
}

func (s *skAppMetricMiddleware) HealthCheck() bool {
	panic("implement me")
}

func (s *skAppMetricMiddleware) SecInfo(productId int) (data map[string]interface{}) {
	panic("implement me")
}

func (s *skAppMetricMiddleware) SecKill(req *model.SecRequest) (map[string]interface{}, int, error) {
	panic("implement me")
}

func (s *skAppMetricMiddleware) SecInfoList() ([]map[string]interface{}, int, error) {
	panic("implement me")
}

func NewSkAppMetrics(requestCount metrics.Counter, requestLatency metrics.Histogram) service.ServiceMiddleware {
	return func(next service.Service) service.Service {
		return skAppMetricMiddleware{
			service:        next,
			requestCount:   requestCount,
			requestLatency: requestLatency,
		}
	}
}
