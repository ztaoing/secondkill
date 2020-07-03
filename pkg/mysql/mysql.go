/**
* @Author:zhoutao
* @Date:2020/7/3 上午8:15
 */

package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose/v2"
)

var (
	err   error
	engin *gorose.Engin
)

func InitMysql(hostMysql, portMysql, userMysql, pwdMysql, dbMysql string) {

	fmt.Printf(userMysql)
	fmt.Printf(dbMysql)

	DBConfig := gorose.Config{
		Driver: "mysql",
		Dsn:    userMysql + ":" + pwdMysql + "@tcp(" + hostMysql + ":" + portMysql + ")/" + dbMysql + "?charset=utf8&parseTime=true", //数据库连接
		Prefix: "",
		//最大连接池
		SetMaxOpenConns: 300,
		//最大空闲连接
		SetMaxIdleConns: 10,
	}

	engin, err = gorose.Open(&DBConfig)
	if err != nil {
		fmt.Println(err)
		return
	}

}

func DB() gorose.IOrm {
	return engin.NewOrm()
}
