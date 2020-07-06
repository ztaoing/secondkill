/**
* @Author:zhoutao
* @Date:2020/7/6 上午10:58
 */

package plugins

import (
	"github.com/go-kit/kit/log"
	"secondkill/sk-app/model"
	"secondkill/sk-app/service"
	"time"
)

//日志中间件
type skAppLoggingMiddleware struct {
	service service.Service
	logger  log.Logger
}

//健康检查
func (s *skAppLoggingMiddleware) HealthCheck() (result bool) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"function", "HealthCheck",
			"result", result,
			"took", time.Since(begin), //用时
		)
	}(time.Now())

	result = s.service.HealthCheck()
	return result
}

//秒杀详情
func (s *skAppLoggingMiddleware) SecInfo(productId int) (data map[string]interface{}) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"function", "SecInfo",
			"took", time.Since(begin),
		)
	}(time.Now())
	data = s.service.SecInfo(productId)
	return data
}

//秒杀
func (s *skAppLoggingMiddleware) SecKill(req *model.SecRequest) (map[string]interface{}, int, error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"function", "SecKill",
			"took", time.Since(begin),
		)
	}(time.Now())
	result, code, err := s.service.SecKill(req)
	return result, code, err
}

//秒杀列表
func (s *skAppLoggingMiddleware) SecInfoList() ([]map[string]interface{}, int, error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"function", "SecKill",
			"took", time.Since(begin),
		)
	}(time.Now())
	result, code, err := s.service.SecInfoList()
	return result, code, err
}

func NewSkAppLoggingMiddleware(logger log.Logger) service.ServiceMiddleware {
	return func(next service.Service) service.Service {
		return &skAppLoggingMiddleware{
			service: next,
			logger:  logger,
		}
	}
}
