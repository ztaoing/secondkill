/**
* @Author:zhoutao
* @Date:2020/7/3 上午8:57
 */

package model

type UserDetails struct {
	//用户标识
	UserId int64
	//用户名 唯一
	UserName string
	//密码
	Password string
	//用户拥有的权限
	Authorities []string
}

func (userDetails *UserDetails) IsMatch(username string, password string) bool {
	return userDetails.Password == password && userDetails.UserName == username
}
