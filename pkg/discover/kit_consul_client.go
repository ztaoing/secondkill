/**
* @Author:zhoutao
* @Date:2020/6/30 下午4:19
* 基于consul的服务发现和注册
 */

package discover

import (
	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"log"
	"secondkill/common"
	"strconv"
	"sync"
)

func New(consulHost, consulPort string) *DiscoverClientInstance {
	port, _ := strconv.Atoi(consulPort)
	//通过consul host 和 consul port 创建一个consul.client
	consulConfig := api.DefaultConfig()
	consulConfig.Address = consulHost + ":" + consulPort

	apiclient, err := api.NewClient(consulConfig)
	if err != nil {
		return nil
	}
	client := consul.NewClient(apiclient)
	return &DiscoverClientInstance{
		Host:   consulHost,
		Port:   port,
		config: consulConfig,
		client: client,
	}
}

//实现DiscoveryClient接口
type DiscoverClientInstance struct {
	Host   string
	Port   int
	config *api.Config
	client consul.Client
	mutex  sync.Mutex
	//服务实例缓存
	instanceMap sync.Map
}

//注册服务
func (consulClient *DiscoverClientInstance) Register(instanceId, svcHost, svcPort, healthCheckUrl string, svcName string, weight int, meta map[string]string, tags []string, logger *log.Logger) bool {
	port, _ := strconv.Atoi(svcPort)

	//构建服务实例元数据
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
	err := consulClient.client.Register(serviceRegistration)
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

//服务注销
func (consulClient *DiscoverClientInstance) DeRegister(instanceId string, logger *log.Logger) bool {
	//构建服务实例元数据
	serviceRegistration := &api.AgentServiceRegistration{
		ID: instanceId,
	}
	//发送服务到consul
	err := consulClient.client.Deregister(serviceRegistration)
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

//查询服务实例
func (consulClient *DiscoverClientInstance) DiscoverServices(serviceName string, logger *log.Logger) []*common.ServiceInstance {
	//是否已经缓存
	instanceList, ok := consulClient.instanceMap.Load(serviceName)
	if ok {
		return instanceList.([]*common.ServiceInstance)
	}

	//缓存中没有
	consulClient.mutex.Lock()
	//再次检查是否以监控
	instanceList, ok = consulClient.instanceMap.Load(serviceName)
	if ok {
		return instanceList.([]*common.ServiceInstance)
	} else {
		//注册并监控
		//启动goroutine对consul上的服务进行监控
		go func() {
			params := make(map[string]interface{})
			//定义param
			params["type"] = "service"
			params["service"] = serviceName
			plan, _ := watch.Parse(params)
			plan.Handler = func(u uint64, i interface{}) {
				if i == nil {
					return
				}
				v, ok := i.([]*api.ServiceEntry)
				if !ok {
					return //数据异常
				}
				if len(v) == 0 {
					//没有服务实例在线
					consulClient.instanceMap.Store(serviceName, []*common.ServiceInstance{})
				} else {
					var healthServices []*common.ServiceInstance
					for _, service := range v {
						//只获取健康的服务实例
						if service.Checks.AggregatedStatus() == api.HealthPassing {
							healthServices = append(healthServices, newServiceInstances(service.Service))
						}
					}
					consulClient.instanceMap.Store(serviceName, healthServices)
				}
			}

			defer plan.Stop()
			plan.Run(consulClient.config.Address)
		}()
	}
	defer consulClient.mutex.Unlock()

	//根据服务名称请求实例列表
	entries, _, err := consulClient.client.Service(serviceName, "", false, nil)
	if err != nil {
		consulClient.instanceMap.Store(serviceName, []*common.ServiceInstance{})
		if logger != nil {
			logger.Println("discover service error")
		}
	}
	instances := make([]*common.ServiceInstance, len(entries))
	for i := 0; i < len(instances); i++ {
		instances[i] = newServiceInstances(entries[i].Service)
	}
	//更新到缓存中
	consulClient.instanceMap.Store(serviceName, instances)
	return instances
}

//解析实例对象
func newServiceInstances(service *api.AgentService) *common.ServiceInstance {
	//如果没有设置 rpcPort
	rpcPort := service.Port - 1
	//如果已经设置rpcPort
	if service.Meta != nil {
		if rpcPortString, ok := service.Meta["rpcPort"]; ok {
			rpcPort, _ = strconv.Atoi(rpcPortString)
		}
	}
	return &common.ServiceInstance{
		Host:     service.Address,
		Port:     service.Port,
		GrpcPort: rpcPort,
		Weight:   service.Weights.Passing,
	}

}
