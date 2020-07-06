/**
* @Author:zhoutao
* @Date:2020/7/6 下午11:17
 */

package main

import "secondkill/sk-core/setup"

func main() {
	setup.InitZk()
	setup.InitRedis()
	setup.RunService()
}
