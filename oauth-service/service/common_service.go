/**
* @Author:zhoutao
* @Date:2020/7/3 上午9:48
 */

package service

type Service interface {
	HealthCheck() bool
}

type CommonService struct {
}

//健康检查，这里仅仅返回true
func (c *CommonService) HealthCheck() bool {
	return true
}

func NewCommonService() *CommonService {
	return &CommonService{}
}
