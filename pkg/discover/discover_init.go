/**
* @Author:zhoutao
* @Date:2020/7/1 上午7:42
 */

package discover

import (
	"errors"
	"log"
	"os"
	"secondkill/common"
	"secondkill/pkg/bootstrap"
	"secondkill/pkg/loadbalance"
)

var ConsulService DiscoveryClient
var LoadBalance loadbalance.LoadBalance
var Logger *log.Logger
var NoInstanceExistedErr = errors.New("no available client")

func init() {
	//实例化一个consul客户端，此处实例化了原生态实现版本
	ConsulService = New(bootstrap.DiscoverConfig.Host, bootstrap.DiscoverConfig.Port)
	LoadBalance = new(loadbalance.RandomLoadBalancce)
	Logger = log.New(os.Stderr, "", log.LstdFlags)
}

//服务发现
func DiscoveryService(serviceName string) (*common.ServiceInstance, error) {
	instances := ConsulService.DiscoverServices(serviceName, Logger)

	if len(instances) < 1 {
		Logger.Println("no available client for :%s", serviceName)
		return nil, NoInstanceExistedErr
	}
	//负载均衡
	return LoadBalance.SelectService(instances)
}
