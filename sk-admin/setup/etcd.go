/**
* @Author:zhoutao
* @Date:2020/7/7 下午10:21
 */

package setup

import (
	"github.com/coreos/etcd/clientv3"
	"log"
	pkgConfig "secondkill/pkg/config"
	"time"
)

func InitEtcd() {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		log.Printf("Connect to etcd failed. ERROR:%v", err)
		return
	}
	//赋值
	pkgConfig.Etcd.EtcdSecProductKey = "product"
	pkgConfig.Etcd.EtcdConn = client

}
