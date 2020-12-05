/**
* @Author:zhoutao
* @Date:2020/7/7 下午10:26
 */

package main

import (
	"secondkill/pkg/bootstrap"
	pkgConfig "secondkill/pkg/config"
	"secondkill/pkg/mysql"
	"secondkill/sk-admin/setup"
)

func main() {
	mysql.InitMysql(pkgConfig.MysqlConfig.Host, pkgConfig.MysqlConfig.Port, pkgConfig.MysqlConfig.User, pkgConfig.MysqlConfig.Pwd, pkgConfig.MysqlConfig.Db)
	setup.InitZK()
	setup.InitHTTP(bootstrap.HttpConfig.Host, bootstrap.HttpConfig.Port)
}
