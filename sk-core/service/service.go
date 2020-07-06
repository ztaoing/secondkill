/**
* @Author:zhoutao
* @Date:2020/7/6 下午10:04
 */

package service

type Service interface {
	SecKill() (int, error)
}

type SecKillService struct {
}

//todo
func (s SecKillService) SecKill(a, b int) (int, error) {
	return (a + b), nil
}
