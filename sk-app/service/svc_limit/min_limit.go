/**
* @Author:zhoutao
* @Date:2020/7/5 下午2:54
* 分钟限制
 */

package svc_limit

/**
分钟限制
*/

type TimeLimit interface {
	Count(nowTime int64) (curCount int)
	Check(nowTime int64) int
}

//分钟限制
type MinLimit struct {
	count   int
	CurTime int64
}

//在1分钟之内访问的次数
func (p *MinLimit) Count(nowTime int64) (curCount int) {
	//大于一分钟
	if nowTime-p.CurTime > 60 {
		p.count = 1
		p.CurTime = nowTime
		curCount = p.count
		return
	}
	//一分钟内
	p.count++
	curCount = p.count
	return
}

//检查用户的访问次数
func (p *MinLimit) Check(nowTime int64) int {
	//一分钟之外
	if nowTime-p.CurTime > 60 {
		return 0
	}
	//一分钟之内
	return p.count
}
