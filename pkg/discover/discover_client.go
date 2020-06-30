package discover

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	"secondkill/common"
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
