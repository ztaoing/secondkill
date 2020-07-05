/**
* @Author:zhoutao
* @Date:2020/7/5 下午3:17
 */

package svc_limit

import (
	"fmt"
	"log"
	"secondkill/pkg/config"
	"secondkill/sk-app/model"
	"sync"
)

type Limit struct {
	secLimit TimeLimit //秒限制
	minLimit TimeLimit //分限制
}

//限制管理
type SecLimitMgr struct {
	UserLimitMap map[int]*Limit
	IpLimitMap   map[string]*Limit
	lock         sync.Mutex
}

var SecLimitMgrVars = &SecLimitMgr{
	UserLimitMap: make(map[int]*Limit),
	IpLimitMap:   make(map[string]*Limit),
}

//黑名单+分、秒限制
func AntiSpam(req *model.SecRequest) (err error) {
	//判断用户ID是否在黑名单中
	_, ok := config.SecKill.IDBlackMap[req.UserId]
	if ok {
		err = fmt.Errorf("invalid request")
		log.Printf("user[%v] is blocked by IDBlack List", req.UserId)
		return
	}
	//是否在客户端ip黑名单中
	_, ok = config.SecKill.IPBlackMap[req.ClientAddr]
	if ok {
		err = fmt.Errorf("invalid request")
		log.Printf("user[%v] ip[%v] is blocked by IPBlack List", req.UserId, req.ClientAddr)
		return
	}

	var secIdCount, minIdCount, secIpCount, minIpCount int
	//加锁
	SecLimitMgrVars.lock.Lock()
	{
		//用户ID频率控制
		limit, ok := SecLimitMgrVars.UserLimitMap[req.UserId]
		if !ok {
			//如果没有就创建
			limit = &Limit{
				secLimit: &SecLimit{},
				minLimit: &MinLimit{},
			}
			SecLimitMgrVars.UserLimitMap[req.UserId] = limit
		}
		//有
		//该秒内用户访问次数
		secIdCount = limit.secLimit.Count(req.AccessTime)
		//该分钟内用户访问次数
		minIdCount = limit.minLimit.Count(req.AccessTime)

		//客户端ip频率控制
		limit, ok = SecLimitMgrVars.IpLimitMap[req.ClientAddr]
		if !ok {
			limit = &Limit{
				secLimit: &SecLimit{},
				minLimit: &MinLimit{},
			}
			SecLimitMgrVars.IpLimitMap[req.ClientAddr] = limit
		}
		//该秒内访问次数
		secIpCount = limit.secLimit.Count(req.AccessTime)
		minIpCount = limit.minLimit.Count(req.AccessTime)

	}
	//释放锁
	SecLimitMgrVars.lock.Unlock()

	//该用户一秒内访问次数 大于 最大访问次数
	if secIdCount > config.SecKill.AccessLimitConf.UserSecAccessLimit {
		err = fmt.Errorf("invalid request")
		return
	}

	//在用户一分钟内访问次数 大于 最大访问次数
	if minIdCount > config.SecKill.AccessLimitConf.UserMinAccessLimit {
		err = fmt.Errorf("invalid request")
		return
	}

	//该IP一秒内访问次数 大于 最大访问次数
	if secIpCount > config.SecKill.AccessLimitConf.IPSecAccessLimit {
		err = fmt.Errorf("invalid request")
		return
	}

	//该IP一分钟内访问次数 大于最大访问次数
	if minIdCount > config.SecKill.AccessLimitConf.IPMinAccessLimit {
		err = fmt.Errorf("invalid request")
		return
	}

	//通过 用户和ip、黑名单校验
	return
}
