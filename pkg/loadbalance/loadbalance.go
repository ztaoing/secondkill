/**
* @Author:zhoutao
* @Date:2020/6/30 下午5:54
 */

package loadbalance

import (
	"errors"
	"math/rand"
	"secondkill/common"
)

type LoadBalance interface {
	//根据负载均衡策略从实例列表中选择一个实例
	SelectService(serviceList []*common.ServiceInstance) (*common.ServiceInstance, error)
}

//完全随机策略
type RandomLoadBalancce struct {
}

func (rl *RandomLoadBalancce) SelectService(serviceList []*common.ServiceInstance) (*common.ServiceInstance, error) {
	if serviceList == nil || len(serviceList) == 0 {
		return nil, errors.New("service instances are not exist")
	}
	return serviceList[rand.Intn(len(serviceList))], nil
}

//带权重的平滑负载策略
type WeightLoadBalance struct {
}

/*
带权重的平滑策略负载均衡需要根据serviceinstance 中的weight 和Curweight这两个属性进行计算
weight配置的服务实例权重固定不变
curweight 是服务实例目前的权重，一开始为0，之后会动态调整
每次当请求到来，选取服务实例时，该策略会遍历服务实例队列中的所有服务实例。对于每个服务实例，让他的CurWeight增加他的weight值；同时累加所有服务实例的weight的和，保存为total
遍历完所有服务实例之后，如果某个服务实例的CurWeight最大，就选择这个服务实例处理本次请求，最后把把该服务实例的curweight减去total值
*/
func (wl *WeightLoadBalance) SelectService(serviceList []*common.ServiceInstance) (best *common.ServiceInstance, err error) {
	if serviceList == nil || len(serviceList) == 0 {
		return nil, errors.New("service instances do not exit")
	}

	total := 0

	for i := 0; i < len(serviceList); i++ {
		w := serviceList[i]
		if w == nil {
			continue
		}

		w.CurWeight += w.Weight
		//统计所有权重之和
		total += w.Weight
		//恢复权重
		if w.CurWeight < w.Weight {
			w.Weight++
		}
		//选择最大临时权重节点
		if best == nil || w.CurWeight > best.CurWeight {
			best = w
		}
	}
	if best == nil {
		return nil, nil
	}
	best.CurWeight -= total
	return best, nil
}
