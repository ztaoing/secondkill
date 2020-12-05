/**
* @Author:zhoutao
* @Date:2020/7/6 上午11:14
 */

package plugins

import (
	"github.com/go-kit/kit/metrics"
	"secondkill/sk-app/model"
	"secondkill/sk-app/service"
	"time"
)

//metrics中间件
type skAppMetricMiddleware struct {
	service                  service.Service
	requestCount             metrics.Counter
	requestLaDockerfiletency metrics.Histogram
}

func (s *skAppMetricMiddleware) HealthCheck() bool {
	defer func(begin time.Time) {
		lvs := []string{
			"method", "HealthCheck",
		}
		s.requestCount.With(lvs...).Add(1)
		s.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	result := s.service.HealthCheck()
	return result
}

func (s *skAppMetricMiddleware) SecInfo(productId int) (data map[string]interface{}) {
	defer func(begin time.Time) {
		lvs := []string{
			"method", "SecInfo",
		}
		s.requestCount.With(lvs...).Add(1)
		s.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	data = s.service.SecInfo(productId)
	return data

}

func (s *skAppMetricMiddleware) SecKill(req *model.SecRequest) (map[string]interface{}, int, error) {
	defer func(begin time.Time) {
		lvs := []string{
			"method", "SecKill",
		}
		s.requestCount.With(lvs...).Add(1)
		s.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	result, code, err := s.service.SecKill(req)
	return result, code, err
}

func (s *skAppMetricMiddleware) SecInfoList() ([]map[string]interface{}, int, error) {
	defer func(begin time.Time) {
		lvs := []string{
			"method", "SecKill",
		}
		s.requestCount.With(lvs...).Add(1)
		s.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	data, code, err := s.service.SecInfoList()
	return data, code, err
}

func NewSkAppMetricsMiddleware(requestCount metrics.Counter, requestLatency metrics.Histogram) service.ServiceMiddleware {
	return func(next service.Service) service.Service {
		return &skAppMetricMiddleware{
			service:        next,
			requestCount:   requestCount,
			requestLatency: requestLatency,
		}
	}
}
