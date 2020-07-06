/**
* @Author:zhoutao
* @Date:2020/7/6 下午3:30
 */

package config

type SecRequest struct {
	ProductId int    `json:"product_id"`
	UserId    int    `json:"user_id"`
	Token     string `json:"token"`
	TokenTime int64  `json:"token_time"`
	Code      int    `json:"code"`
}
type SecResult struct {
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
