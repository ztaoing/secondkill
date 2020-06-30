/**
* @Author:zhoutao
* @Date:2020/6/30 下午4:19
 */

package discover

import (
	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	"strconv"
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
