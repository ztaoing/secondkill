/**
* @Author:zhoutao
* @Date:2020/7/6 上午7:42
 */

package setup

import (
	"github.com/go-redis/redis"
	"github.com/unknwon/com"
	"log"
	"secondkill/pkg/config"
	"secondkill/sk-app/service/svc_redis"
	"time"
)

//初始化redis
func InitRedis() {
	log.Printf("init redis")

	client := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Host,
		Password: config.Redis.Password,
		DB:       config.Redis.Db,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Printf("Connect redis failed . ERROR:%v", err)
		return
	}
	log.Printf("init redis success")
	config.Redis.RedisConn = client
	//加载黑名单信息
	loadBlackList(client)

	initRedisProcess()
}

//加载黑名单
func loadBlackList(conn *redis.Client) {
	config.SecKill.IPBlackMap = make(map[string]bool, 10000)
	config.SecKill.IDBlackMap = make(map[int]bool, 10000)

	//用户ID限制
	idList, err := conn.HGetAll(config.Redis.IdBlackListHash).Result()
	if err != nil {
		log.Printf("hget_all IdBlackListHash failed. ERROR:%v", err)
		return
	}
	for _, v := range idList {
		id, err := com.StrTo(v).Int()
		if err != nil {
			log.Printf("invalid user id:%v", id)
			continue
		}
		//放到本地缓存中
		config.SecKill.IDBlackMap[id] = true
	}

	//用户ip限制
	ipList, err := conn.HGetAll(config.Redis.IpBlackListHash).Result()
	if err != nil {
		log.Printf("hget_all IpBlackListHash failed:ERROR:%v", err)
		return
	}
	for _, v := range ipList {
		config.SecKill.IPBlackMap[v] = true
	}
	go syncIpBlackList(conn)
	go syncIdBlackList(conn)
	return
}

//初始化redis
func initRedisProcess() {
	log.Printf("initRedisProcess write:%d,read: %d", config.SecKill.AppWriteToHandleGoroutineNum, config.SecKill.AppReadFromHandleGoroutineNum)
	for i := 0; i < config.SecKill.AppWriteToHandleGoroutineNum; i++ {
		go svc_redis.WriteHandler()
	}

	for i := 0; i < config.SecKill.AppReadFromHandleGoroutineNum; i++ {
		go svc_redis.ReadHandler()
	}
}

//同步用户id黑名单
func syncIdBlackList(conn *redis.Client) {
	for {
		//阻塞获取队列中的数据
		idArr, err := conn.BRPop(time.Minute, config.Redis.IdBlackListQueue).Result()
		if err != nil {
			//阻塞式取出
			log.Printf("brpop id failed, ERROR:%v", err)
			continue
		}
		id, _ := com.StrTo(idArr[1]).Int()
		config.SecKill.RWBlackLock.Lock()

		config.SecKill.IDBlackMap[id] = true

		config.SecKill.RWBlackLock.Unlock()

	}
}

//同步用户ip黑名单
func syncIpBlackList(conn *redis.Client) {
	var ipList []string
	lastTime := time.Now().Unix()
	for {
		ipArr, err := conn.BRPop(time.Minute, config.Redis.IpBlackListQueue).Result()
		if err != nil {
			log.Printf("brpop is faield,ERROR:%v", err)
			continue
		}
		ip := ipArr[1]
		curTime := time.Now().Unix()
		ipList = append(ipList, ip)
		//数量+时间限制
		if len(ipList) > 100 || curTime-lastTime > 5 {
			config.SecKill.RWBlackLock.Lock()
			{
				for _, v := range ipList {
					config.SecKill.IPBlackMap[v] = true
				}
			}
			config.SecKill.RWBlackLock.Unlock()

			lastTime = curTime
			log.Printf("sync ip list form redis success,IP[%v]", ipList)
		}

	}
}
