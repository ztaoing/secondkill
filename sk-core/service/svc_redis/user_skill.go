/**
* @Author:zhoutao
* @Date:2020/7/6 下午6:20
 */

package svc_redis

import (
	"crypto/md5"
	"fmt"
	"log"
	pkgConfig "secondkill/pkg/config"
	"secondkill/sk-core/config"
	"secondkill/sk-core/service/svc_err"
	"secondkill/sk-core/service/svc_user"
	"time"
)

//HandleUser 从Read2HandleChan 中获取请求，然后调用HandleSecKill最用户秒杀请求进行处理，
//将返回结果推入write2HandleChan 中，等待结果写入Redis，并设置结果写入redis操作的超时时间和超时回调
func HandleUser() {
	log.Println("handle user running")
	for req := range config.SecLayerCtx.Read2HandleChan {
		log.Printf("begin process request:%v", req)
		//处理抢购
		res, err := HandleSecKill(req)
		if err != nil {
			log.Printf("process request [%v] failed,ERROR:%v", req, err)
			res = &config.SecResult{
				Code: svc_err.ErrServiceBusy,
			}
		}

		timer := time.NewTicker(time.Millisecond * time.Duration(pkgConfig.SecKill.SendToHandChanTimeout))
		select {
		case config.SecLayerCtx.Write2HandleChan <- res:
		case <-timer.C:
			log.Printf("send to response chan timeout,res:%v", res)
			break
		}
	}
	return
}

//限制用户对商品的购买次数，对商品的抢购概率进行限制。对合法的请求生成抢购资格token令牌
func HandleSecKill(req *config.SecRequest) (res *config.SecResult, err error) {
	config.SecLayerCtx.RWSecProductLock.RLock()
	defer config.SecLayerCtx.RWSecProductLock.RUnlock()

	res = &config.SecResult{}
	res.ProductId = req.ProductId
	res.UserId = req.UserId

	//获得商品详情
	product, ok := pkgConfig.SecKill.SecProductInfoMap[req.ProductId]
	if !ok {
		log.Printf("not found product:%v", req.ProductId)
		res.Code = svc_err.ErrNotFoundProduct
		return
	}

	if product.Status == svc_err.ProductStatusSoldOut {
		res.Code = svc_err.ErrSoldOut
		return
	}

	nowTime := time.Now().Unix()

	//加锁 是否已经购买过
	config.SecLayerCtx.HistoryMapLock.Lock()
	userHistory, ok := config.SecLayerCtx.HistoryMap[req.UserId]

	if !ok {
		userHistory = &svc_user.UserBuyHistory{
			History: make(map[int]int, 16),
		}
		config.SecLayerCtx.HistoryMap[req.UserId] = userHistory
	}

	historyCount := userHistory.GetProductBuyCount(req.ProductId)
	config.SecLayerCtx.HistoryMapLock.Unlock()

	//大于个人购买商品的数量
	if historyCount >= product.OnePersonBuyLimit {
		res.Code = svc_err.ErrAlreadyBuy
	}
	curSoldCount := config.SecLayerCtx.ProductCountMgr.Count(req.ProductId)
	//超卖
	if curSoldCount >= product.Total {
		res.Code = svc_err.ErrSoldOut
		product.Status = svc_err.ProductStatusSoldOut
		return
	}

	curRate := 0.1
	fmt.Println(curRate, product.BuyRate)

	if curRate > product.BuyRate {
		res.Code = svc_err.ErrRetry
		return
	}

	//通过以上限制
	userHistory.Add(req.ProductId, 1)
	//增加数量，add内已加锁
	config.SecLayerCtx.ProductCountMgr.Add(req.ProductId, 1)

	//
	res.Code = svc_err.SecKillSuccess
	tokenData := fmt.Sprintf("userId=%d&productId=%d&timestamp=%d&security=%s",
		req.UserId,
		req.ProductId,
		nowTime,
		pkgConfig.SecKill.TokenPassWd,
	)
	res.Token = fmt.Sprintf("%x", md5.Sum([]byte(tokenData)))
	res.TokenTime = nowTime
	return
}
