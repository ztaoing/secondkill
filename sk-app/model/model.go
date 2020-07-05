/**
* @Author:zhoutao
* @Date:2020/7/5 下午2:09
 */

package model

//请求
type SecRequest struct {
	ProductId       int             `json:"product_id"`
	Source          string          `json:"source"`
	AuthCode        string          `json:"auth_code"`
	SecTime         int64           `json:"sec_time"`
	Nance           string          `json:"nance"`
	UserId          int             `json:"user_id"`
	UserAuthSign    string          `json:"user_auth_sign"`
	AccessTime      int64           `json:"access_time"`
	ClientAddr      string          `json:"client_addr"`
	ClientReference string          `json:"client_reference"`
	CloseNotify     <-chan bool     `json:"-"`
	ResultChan      chan *SecResult `json :"-"'`
}

type SecResult struct {
	ProductId int    `json:"product_id"`
	UserId    int    `json:"user_id"`
	Token     string `json:"token"`
	TokenTime int64  `json:"token_time"` //token生成的时间
	Code      int    `json:"code"`       //状态码
}
