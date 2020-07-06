/**
* @Author:zhoutao
* @Date:2020/7/6 下午3:30
 */

package config

import (
	"secondkill/sk-core/service/svc_product"
	"secondkill/sk-core/service/svc_user"
	"sync"
)

type SecResult struct {
	ProductId int    `json:"product_id"`
	UserId    int    `json:"user_id"`
	Token     string `json:"token"`
	TokenTime int64  `json:"token_time"`
	Code      int    `json:"code"`
}
type SecRequest struct {
	ProductId       int             `json:"product_id"`
	Source          string          `json:"source"`
	AuthCode        string          `json:"auth_code"`
	SecTime         int64           `json:"sec_time"`
	Nance           string          `json:"nance"`
	UserId          int             `json:"user_id"`
	UserAuthSign    string          `json:"user_auth_sign"`
	ClientAddr      string          `json:"client_addr"`
	ClientReference string          `json:"client_reference"`
	CloseNotify     <-chan bool     `json:"-"`
	ResultChan      chan *SecResult `json:"-"`
}

type SecLayerContext struct {
	RWSecProductLock sync.RWMutex

	WaitGroup sync.WaitGroup

	Read2HandleChan  chan *SecRequest
	Write2HandleChan chan *SecResult

	HistoryMap     map[int]*svc_user.UserBuyHistory
	HistoryMapLock sync.Mutex

	ProductCountMgr *svc_product.ProductCountMgr
}

var SecLayerCtx = &SecLayerContext{
	Read2HandleChan:  make(chan *SecRequest, 1024),
	Write2HandleChan: make(chan *SecResult, 1024),
	HistoryMap:       make(map[int]*svc_user.UserBuyHistory, 1024),
	ProductCountMgr:  svc_product.NewProductCountMgr(),
}
