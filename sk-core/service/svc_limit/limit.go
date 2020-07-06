/**
* @Author:zhoutao
* @Date:2020/7/6 下午4:24
 */

package svc_limit

//每秒访问限制
type SecLimit struct {
	count   int   //秒内访问次数
	preTime int64 //上一次记录的时间
}

func (p *SecLimit) Count(nowTime int64) (curCount int) {
	if p.preTime != nowTime {
		p.count = 1
		p.preTime = nowTime
		curCount = p.count
		return
	}
	p.count++
	curCount = p.count
	return
}

func (p *SecLimit) Check(nowTime int64) int {
	if p.preTime != nowTime {
		return 0
	}
	return p.count
}
