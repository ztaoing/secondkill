/**
* @Author:zhoutao
* @Date:2020/7/7 上午10:05
 */

package service

type Service interface {
	HealthCheck() bool
}

type SKAdminService struct {
}

func (s *SKAdminService) HealthCheck() bool {
	return true
}

//装饰者
type ServiceMiddleware func(Service) Service
