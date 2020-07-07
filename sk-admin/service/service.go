/**
* @Author:zhoutao
* @Date:2020/7/7 上午10:05
 */

package service

type service interface {
	HealthCheck() bool
}

type SKAdminService struct {
}

func (s *SKAdminService) HealthCheck() bool {
	return true
}

type ServiceMiddleware func(service) service
