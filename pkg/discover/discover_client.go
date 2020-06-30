package discover

import (
	"fmt"
	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	"log"
	"secondkill/common"
	"strconv"
	"sync"
)

type DiscoveryClient interface {
	/**
	服务注册接口
	@param instanceId 实例ID
	@param svcHost 服务的host
	@param svcPort 服务的port
	@param healthCheckUrl 健康检查的地址
	@param svcName 服务的名称
	@param weight 权重
	@param meta map[string]string 服务实例元数据
	@param tags []string
	*/
	Register(instanceId, svcHost, svcPort, healthCheckUrl string, svcName string, weight int, meta map[string]string, tags []string, logger *log.Logger) bool

	/**
	服务注销接口
	@param instanceId 实例的id
	*/
	DeRegister(instanceId, logger *log.Logger) bool

	/**
	服务发现
	@param serviceName 服务的名称
	*/
	DiscoverServices(serviceName string, logger *log.Logger) []*common.ServiceInstance
}

//会实现DiscoveryClient接口
type DiscoverClientInstance struct {
	Host   string
	Port   int
	config *api.Config
	cleint consul.Client
	mutex  sync.Mutex
	//服务实例缓存
	instanceMap sync.Map
}

func (consulClient *DiscoverClientInstance) Register(instanceId, svcHost, svcPort, healthCheckUrl string, svcName string, weight int, meta map[string]string, tags []string, logger *log.Logger) bool {
	port, _ := strconv.Atoi(svcPort)

	//构建服务实例元数据
	fmt.Println(weight)
	serviceRegistration := &api.AgentServiceRegistration{
		ID:      instanceId,
		Name:    svcName,
		Address: svcHost,
		Port:    port,
		Meta:    meta,
		Tags:    tags,
		Weights: &api.AgentWeights{
			Passing: weight,
		},
		Check: &api.AgentServiceCheck{
			DeregisterCriticalServiceAfter: "30s",
			HTTP:                           "http://" + svcHost + ":" + svcPort + healthCheckUrl,
			Interval:                       "15s",
		},
	}

	//发送服务到consul
	err := consulClient.cleint.Register(serviceRegistration)
	if err != nil {
		if logger != nil {
			logger.Println("Register service error")
		}
		return false
	}
	if logger != nil {
		logger.Println("Register service success")
	}
	return true

}

func (consulClient *DiscoverClientInstance) DeRegister(instanceId string, logger *log.Logger) bool {
	//构建服务实例元数据
	serviceRegistration := &api.AgentServiceRegistration{
		ID: instanceId,
	}
	//发送服务到consul
	err := consulClient.cleint.Deregister(serviceRegistration)
	if err != nil {
		if logger != nil {
			logger.Println("DeRegister service error")
		}
		return false
	}
	if logger != nil {
		logger.Println("DeRegister service success")
	}
	return true
}

func (consulClient *DiscoverClientInstance) DiscoverServices(serviceName string, logger *log.Logger) []*common.ServiceInstance {
	//是否已经缓存

}
