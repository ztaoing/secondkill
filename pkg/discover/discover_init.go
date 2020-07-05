/**
* @Author:zhoutao
* @Date:2020/7/1 上午7:42
 */

package discover

import (
	"errors"
	uuid "github.com/satori/go.uuid"
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

//服务注册
func Register() {
	if ConsulService == nil {
		panic("ConsulService failed")
	}

	instanceId := bootstrap.DiscoverConfig.InstanceID
	if instanceId == "" {
		instanceId = bootstrap.DiscoverConfig.ServiceName + uuid.NewV4().String()
	}
	if !ConsulService.Register(instanceId,
		bootstrap.DiscoverConfig.Host,
		bootstrap.DiscoverConfig.Port,
		"/health",
		bootstrap.DiscoverConfig.ServiceName,
		bootstrap.DiscoverConfig.Weight,
		map[string]string{
			"rpcPort": bootstrap.RpcConfig.Port,
		}, nil, Logger,
	) {
		//注册失败
		Logger.Printf("register service %s failed", bootstrap.DiscoverConfig.ServiceName)
		panic("consulService register failed")
	}

	//注册成功
	Logger.Printf(bootstrap.DiscoverConfig.ServiceName+"-register for service %s success", bootstrap.DiscoverConfig.ServiceName)

}

//服务注销
func Deregister() {
	if ConsulService == nil {
		panic("consul srvice is nil")
	}
	//
	instanceId := bootstrap.DiscoverConfig.InstanceID

	if instanceId == "" {
		instanceId = bootstrap.DiscoverConfig.ServiceName + "-" + uuid.NewV4().String()
	}
	if !ConsulService.DeRegister(instanceId, Logger) {
		Logger.Printf("deregister for service %s failed", bootstrap.DiscoverConfig.ServiceName)
		panic(0)
	}

}
