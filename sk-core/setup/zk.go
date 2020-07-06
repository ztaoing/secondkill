/**
* @Author:zhoutao
* @Date:2020/7/6 下午10:09
 */

package setup

import (
	"encoding/json"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	pkgConfig "secondkill/pkg/config"
	"time"
)

/**
从zookeeper中加载秒杀活动到内存中，监听zookeeper的数据变化，实时更新数据到内存中
*/
func InitZk() {
	var host = []string{"127.0.0.1:2181"}
	cnnOption := zk.WithEventCallback(waitSecProductEvent)
	conn, _, err := zk.Connect(host, time.Second*5, cnnOption)
	if err != nil {
		fmt.Println(err)
		return
	}

	pkgConfig.Zk.ZkConn = conn
	pkgConfig.Zk.SecProductKey = "/product"

	//加载秒杀商品的信息
	go loadSecConf(conn)

}

func waitSecProductEvent(event zk.Event) {
	if event.Path == pkgConfig.Zk.SecProductKey {

	}
}

func loadSecConf(conn *zk.Conn) {
	log.Printf("Connect zookeeper success %s", pkgConfig.Zk.SecProductKey)
	v, _, err := conn.Get(pkgConfig.Zk.SecProductKey)
	if err != nil {
		log.Printf("get product info from zookeeper failed,ERROR:%v", err)
		return
	}

	log.Printf("get product info success")
	var secProductInfo []*pkgConfig.SecProductInfoConf

	err = json.Unmarshal(v, &secProductInfo)
	if err != nil {
		log.Printf("unmarshal product info failed ,ERROR:%v", err)
	}
	//更新秒杀商品信息
	updateSecProductInfo(secProductInfo)

}

func updateSecProductInfo(secProductInfo []*pkgConfig.SecProductInfoConf) {
	temp := make(map[int]*pkgConfig.SecProductInfoConf, 1024)
	for _, v := range secProductInfo {
		temp[v.ProductID] = v
	}
	//加锁
	pkgConfig.SecKill.RWBlackLock.Lock()
	pkgConfig.SecKill.SecProductInfoMap = temp
	pkgConfig.SecKill.RWBlackLock.Unlock()
}
