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
	"secondkill/sk-app/model"
	"time"
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
		conn := pkgConfig.Redis.RedisConn
		//使用阻塞弹出
		data, err := conn.BRPop(time.Second*1, pkgConfig.Redis.Proxy2layerQueueName).Result()
		if err != nil {
			continue
		}
		var result *model.SecResult
		err = json.Unmarshal([]byte(data[1]), &result)
		if err != nil {
			log.Printf("json unmarshal failed . Error : %v", err)
			continue
		}

		userKey := fmt.Sprintf("%d_%d", result.UserId, result.ProductId)
		fmt.Println("userKey:", userKey)

		//加锁
		config.SkAppContext.UserConnMapLock.Lock()

		resultChan, ok := config.SkAppContext.UserConnMap[userKey]
		//解锁
		config.SkAppContext.UserConnMapLock.Unlock()
		if !ok {
			log.Printf("user not found:%v", userKey)
			continue
		}
		log.Printf("request result send to chan")
		//channel验证通过，将结果放入channel
		resultChan <- result
		log.Printf("request result send to chan success,userKey:%v", userKey)
	}
}
