/**
* @Author:zhoutao
* @Date:2020/7/6 下午10:09
 */

package setup

import (
	"github.com/go-redis/redis"
	"log"
	pkgConfig "secondkill/pkg/config"
)

/**
初始化redis
*/

func InitRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:     pkgConfig.Redis.Host,
		Password: pkgConfig.Redis.Password,
		DB:       pkgConfig.Redis.Db,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Printf("Connect redis failed. ERROR:%v", err)
	}
	pkgConfig.Redis.RedisConn = client
}
