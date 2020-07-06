/**
* @Author:zhoutao
* @Date:2020/7/6 上午7:42
 */

package setup

import (
	"encoding/json"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"secondkill/pkg/config"
	"time"
)

//初始化zookeeper
func InitZk() {
	var hosts = []string{"127.0.0.1:2181"}
	options := zk.WithEventCallback(waitSecProductEvent)
	conn, _, err := zk.Connect(hosts, time.Second*5, options)
	if err != nil {
		fmt.Println(err)
		return
	}
	config.Zk.ZkConn = conn
	config.Zk.SecProductKey = "/product"
	loadSecConf(conn)
}

func waitSecProductEvent(event zk.Event) {
	log.Print("===================")
	log.Println("path:", event.Path)
	log.Println("type:", event.Type.String())
	log.Println("state:", event.State.String())
	log.Print("===================")

	if event.Path == config.Zk.SecProductKey {
		//todo
	}
}

//加载秒杀商品信息
func loadSecConf(conn *zk.Conn) {
	log.Printf("Connect zk success %s", config.Zk.SecProductKey)
	v, _, err := conn.Get(config.Zk.SecProductKey)
	if err != nil {
		log.Printf("get product info failed,ERROR:%v", err)
		return

	}

	//获取商品信息成功
	log.Printf("get product info success,ERROR")

	var secProductInfo []*config.SecProductInfoConf
	err = json.Unmarshal(v, &secProductInfo)
	if err != nil {
		log.Printf("json unmarshal secProductInfo failed,ERROR:%v", err)
		return
	}
	updateSecProductInfo(secProductInfo)
}

func updateSecProductInfo(secProductInfo []*config.SecProductInfoConf) {
	temp := make(map[int]*config.SecProductInfoConf, 1024)
	for _, v := range secProductInfo {
		log.Printf("updateSecProductInfo %v", v)
		temp[v.ProductID] = v
	}
	//更新缓存信息
	config.SecKill.RWBlackLock.Lock()
	config.SecKill.SecProductInfoMap = temp
	config.SecKill.RWBlackLock.Unlock()

}
