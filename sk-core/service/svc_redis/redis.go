/**
* @Author:zhoutao
* @Date:2020/7/6 下午4:48
 */

package svc_redis

import (
	"encoding/json"
	"fmt"
	"log"
	pkgConfig "secondkill/pkg/config"
	"secondkill/sk-core/config"
	"time"
)

//runProcess 会从Redis队列中读取用户请求信息，然后调用goroutine执行HandlerReader操作，并将响应信息推入到redis队列中。

func RunProcess() {
	//读
	for i := 0; i < pkgConfig.SecKill.AppReadFromHandleGoroutineNum; i++ {
		config.SecLayerCtx.WaitGroup.Add(1)
		go HandleReader()
	}

	//写响应
	for i := 0; i < pkgConfig.SecKill.AppWriteToHandleGoroutineNum; i++ {
		config.SecLayerCtx.WaitGroup.Add(1)
		go HandleWriter()
	}

	for i := 0; i < pkgConfig.SecKill.CoreHandleGoroutineNum; i++ {
		config.SecLayerCtx.WaitGroup.Add(1)
		go HandleUser()
	}

	log.Printf("all process goroutine started")
	config.SecLayerCtx.WaitGroup.Wait()
	log.Printf("wait all goroutine exited")

	return
}

func HandleReader() {
	log.Printf("read goroutine [%v] running", pkgConfig.Redis.Proxy2layerQueueName)

	conn := pkgConfig.Redis.RedisConn
	for {
		//从redis队列中读取数据
		data, err := conn.BRPop(time.Second, pkgConfig.Redis.Proxy2layerQueueName).Result()
		if err != nil {
			continue
		}
		log.Printf("block right pop from proxy layer queue,data :%s\n", data)

		//将读取到的请求转换为secRequest struct
		var req config.SecRequest
		err = json.Unmarshal([]byte(data[1]), &req)
		if err != nil {
			log.Printf("unmarshal to secRequest failed,ERROR:%v", err)
			continue
		}

		//判断是否超时
		nowTime := time.Now().Unix()
		//?todo
		fmt.Println(nowTime, " ", req.SecTime, "", 100)
		if nowTime-req.SecTime >= int64(pkgConfig.SecKill.MaxRequestWaitTimeout) {
			log.Printf("request [%v] is expired", req)
			continue
		}
		//设置超时时间
		timer := time.NewTicker(time.Millisecond * time.Duration(pkgConfig.SecKill.CoreWaitResultTimeOut))
		select {
		case config.SecLayerCtx.Read2HandleChan <- &req:
		case <-timer.C:
			log.Printf("send to handle chan timeout,request [%v]", req)
			break
		}
	}
}

func HandleWriter() {
	for res := range config.SecLayerCtx.Write2HandleChan {
		fmt.Printf("===", res)
		err := sendToRedis(res)
		if err != nil {
			log.Printf("send to redis ERROR :%v,response:%v", err, res)
			continue
		}
	}
}

func sendToRedis(res *config.SecResult) (err error) {
	data, err := json.Marshal(res)
	if err != nil {
		log.Printf("marshal failed, ERROR:%v", err)
		return
	}
	fmt.Printf("推入【%v】redis队列", pkgConfig.Redis.Layer2ProxyQueueName)
	conn := pkgConfig.Redis.RedisConn
	err = conn.LPush(pkgConfig.Redis.Layer2ProxyQueueName, string(data)).Err()
	fmt.Printf("推入【%v】redis队列完成", pkgConfig.Redis.Layer2ProxyQueueName)

	if err != nil {
		log.Printf("rpush layer to proxy redis queue failed,ERROR:%v", err)
		return
	}
	log.Printf("Left push success. data[%v]", string(data))
	return
}
