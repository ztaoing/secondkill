/**
* @Author:zhoutao
* @Date:2020/7/5 下午10:17
 */

package service

import (
	"fmt"
	"log"
	"math/rand"
	pkgConfig "secondkill/pkg/config"
	"secondkill/sk-app/config"
	"secondkill/sk-app/model"
	"secondkill/sk-app/service/svc_err"
	"secondkill/sk-app/service/svc_limit"
	"time"
)

type Service interface {
	HealthCheck() bool
	SecInfo(productId int) (data map[string]interface{})
	SecKill(req *model.SecRequest) (map[string]interface{}, int, error)
	SecInfoList() ([]map[string]interface{}, int, error)
}

type ServiceMiddleware func(Service) Service

type SkAppService struct {
}

func (s SkAppService) HealthCheck() bool {
	//此处未做处理，仅返回true
	return true
}

func (s SkAppService) SecInfo(productId int) (data map[string]interface{}) {
	//加锁
	config.SkAppContext.RWSecProductLock.RLock()
	defer config.SkAppContext.RWSecProductLock.Unlock()

	v, ok := pkgConfig.SecKill.SecProductInfoMap[productId]
	if !ok {
		return nil
	}

	data = make(map[string]interface{})
	data["product_id"] = productId
	data["start_time"] = v.StartTime
	data["end_time"] = v.EndTime
	data["status"] = v.Status
	return data
}

func (s SkAppService) SecKill(req *model.SecRequest) (map[string]interface{}, int, error) {
	//对map加锁处理
	config.SkAppContext.RWSecProductLock.RLock()
	defer config.SkAppContext.RWSecProductLock.RUnlock()

	//错误码
	var code int
	//黑名单及分钟、秒限流处理
	err := svc_limit.AntiSpam(req)
	if err != nil {
		code = svc_err.ErrUserServiceBusy
		log.Printf("userId anti_spam [%d] failed,request[%v]", req.UserId, err)
		return nil, code, err
	}

	//获取秒杀信息
	data, code, err := SecInfoById(req.ProductId)
	if err != nil {
		log.Printf("userId[%d] secInfoById Id failed,req[%v]", req.UserId, req)
		return nil, code, err
	}
	userKey := fmt.Sprintf("%d_%d", req.UserId, req.ProductId)

	ResultChan := make(chan *model.SecResult, 1)

	config.SkAppContext.UserConnMapLock.Lock()
	//resultChan
	config.SkAppContext.UserConnMap[userKey] = ResultChan

	config.SkAppContext.UserConnMapLock.Unlock()

	//将请求送入channel，并推入到redis中
	config.SkAppContext.SecReqChan <- req

	ticker := time.NewTicker(time.Millisecond * time.Duration(pkgConfig.SecKill.AppWaitResultTimeOut))

	defer func() {
		ticker.Stop()
		config.SkAppContext.UserConnMapLock.Lock()
		//释放链接
		delete(config.SkAppContext.UserConnMap, userKey)
		config.SkAppContext.UserConnMapLock.Unlock()
	}()

	select {
	case <-ticker.C:
		//超时
		code = svc_err.ErrProcessTimeout
		err = fmt.Errorf("request timeout")
		return nil, code, err

	case <-req.CloseNotify:
		//客户端已关闭
		code = svc_err.ErrClientClose
		err = fmt.Errorf("client already closed")
		return nil, code, err

	case result := <-ResultChan:
		code = result.Code
		//操作不成功
		if code != 1002 {
			return data, code, svc_err.GetErrMsg(code)
		}
		//操作成功
		log.Printf("secKill success!")
		data["product_id"] = result.ProductId
		data["token"] = result.Token
		data["user_id"] = result.UserId

		return data, code, nil
	}

}

func (s SkAppService) SecInfoList() ([]map[string]interface{}, int, error) {
	config.SkAppContext.RWSecProductLock.RLock()
	defer config.SkAppContext.RWSecProductLock.RUnlock()

	var data []map[string]interface{}

	for _, v := range pkgConfig.SecKill.SecProductInfoMap {
		item, _, err := SecInfoById(v.ProductID)
		if err != nil {
			log.Printf("get sec info,Error:%v", err)
			continue
		}
		data = append(data, item)
	}
	return data, 0, nil
}

func NewSecRequest() *model.SecRequest {
	secRequest := &model.SecRequest{
		ResultChan: make(chan *model.SecResult, 1),
	}
	return secRequest
}

func SecInfoById(productId int) (map[string]interface{}, int, error) {
	//加锁
	config.SkAppContext.RWSecProductLock.RLock()
	defer config.SkAppContext.RWSecProductLock.RUnlock()

	var code int
	v, ok := pkgConfig.SecKill.SecProductInfoMap[productId]

	if !ok {
		return nil, svc_err.ErrNotFoundProductId, fmt.Errorf("not find product id:%d", productId)
	}

	//秒杀活动是否开始
	start := false
	//秒杀活动是否结束
	end := false
	//状态
	status := "success"

	var err error
	//当前时间
	nowTime := time.Now().Unix()

	//秒杀活动没有开始
	if nowTime-v.StartTime < 0 {
		start = false
		end = false
		status = "second kill not start"
		code = svc_err.ErrActiveNotStart
		err = fmt.Errorf(status)
	}

	//秒杀活动已经开始
	if nowTime-v.StartTime > 0 {
		start = true
	}

	//秒杀活动已经结束
	if nowTime-v.EndTime > 0 {
		start = false
		end = true
		status = "second kill is already end"
		code = svc_err.ErrActiveAlreadyEnd
		err = fmt.Errorf(status)
	}

	//商品已经售完
	if v.Status == config.ProductStatusForceSaleOut || v.Status == config.ProductStatusSaleOut {
		start = false
		end = false
		status = "product is sale out"
		code = svc_err.ErrActiveSaleOut
		err = fmt.Errorf(status)
	}
	/**
	允许 大于商品数1.5倍的请求进入 秒杀核心层
	*/
	//TODO curRate 随机？
	curRate := rand.Float64()
	if curRate > v.BuyRate*1.5 {
		start = false
		end = false
		status = "retry"
		code = svc_err.ErrRetry
		err = fmt.Errorf(status)
	}

	//组装数据
	data := map[string]interface{}{
		"product_id": productId,
		"start":      start,
		"end":        end,
		"status":     status,
	}
	return data, code, err
}
