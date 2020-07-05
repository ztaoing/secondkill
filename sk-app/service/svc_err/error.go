/**
* @Author:zhoutao
* @Date:2020/7/5 下午2:18
* 定义错误编号
 */

package svc_err

import "errors"

const (
	ErrServiceBusy      = 1001
	SecKillSuccess      = 1002
	ErrNostFoundProduct = 1003
	ErrSoldOut          = 1004
	ErrRetry            = 1005
	ErrAlreadyBuy       = 1006

	ErrInvalidRequest      = 1101
	ErrNotFoundProductId   = 1102
	ErrUserCheckAuthFailed = 1103
	ErrUserServiceBusy     = 1104
	ErrActiveNotStart      = 1105
	ErrActiveAlreadyEnd    = 1106
	ErrActiveSaleOut       = 1107
	ErrProcessTimeout      = 1108
	ErrClientClose         = 1109
)

var errMsg = map[int]string{
	ErrServiceBusy:       "服务器错误",
	SecKillSuccess:       "抢购成功",
	ErrNotFoundProductId: "没有该商品",
	ErrSoldOut:           "商品售完",
	ErrRetry:             "请重试",
	ErrAlreadyBuy:        "已抢购",
}

func GetErrMsg(code int) error {
	return errors.New(errMsg[code])
}
