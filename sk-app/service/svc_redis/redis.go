/**
* @Author:zhoutao
* @Date:2020/7/5 下午4:30
 */

package svc_redis

import (
	"encoding/json"
	"fmt"
	"log"
	pkgConfig "secondkill/pkg/config"
	"secondkill/sk-app/config"
)

//写数据到redis
func WriteHandler() {
	for {
		fmt.Println("write data to redis")

		//读取request
		req := <-config.SkAppContext.SecReqChan

		fmt.Printf("request accessTime :", req.AccessTime)

		//获取redis conn
		conn := pkgConfig.Redis.RedisConn

		data, err := json.Marshal(req)
		if err != nil {
			log.Printf("json.Marshal request failed. Error:%v, Reqest:%v", err, req)
			continue
		}

		//将请求写入到redis中
		err = conn.LPush(pkgConfig.Redis.Proxy2layerQueueName, string(data)).Err()
		if err != nil {
			log.Printf("lpush request failed. Error:%v, Request:%v", err, req)
			continue
		}

		log.Printf("lpush request success. Request:%v", string(data))

	}
}

//从redis读取数据
func ReadHandler() {
	for {

	}
}
