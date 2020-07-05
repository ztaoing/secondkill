/**
* @Author:zhoutao
* @Date:2020/7/5 下午3:00
 */

package svc_limit

/**
秒限制
*/

type SecLimit struct {
	count   int
	curTime int64
}

//一秒内访问的次数
func (p *SecLimit) Count(nowTime int64) (curCount int) {
	//超过一秒
	if p.curTime != nowTime {
		p.count = 1
		p.curTime = nowTime
		curCount = p.count
		return
	}
	//一秒内
	p.count++
	curCount = p.count
	return
}

//检查用户访问次数
func (p *SecLimit) Check(nowTime int64) int {
	if p.curTime != nowTime {
		return 0
	}
	//一秒内的次数
	return p.count
}
