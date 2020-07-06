/**
* @Author:zhoutao
* @Date:2020/7/6 下午2:34
 */

package main

import (
	"secondkill/pkg/bootstrap"
	"secondkill/sk-app/setup"
)

func main() {
	setup.InitZk()
	setup.InitRedis()
	setup.InitHTTP(bootstrap.HttpConfig.Host, bootstrap.HttpConfig.Port)
}
