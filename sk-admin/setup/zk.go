/**
* @Author:zhoutao
* @Date:2020/7/7 下午9:30
 */

package setup

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	pkgConfig "secondkill/pkg/config"
	"time"
)

//初始化zookeeper
func InitZK() {
	var hosts = []string{"127.0.0.1:2181"}
	conn, _, err := zk.Connect(hosts, time.Second*5)
	if err != nil {
		fmt.Println(err)
		return
	}
	pkgConfig.Zk.ZkConn = conn
	pkgConfig.Zk.SecProductKey = "/product"
}
