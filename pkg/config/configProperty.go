/**
* @Author:zhoutao
* @Date:2020/7/1 下午1:54
 */

package config

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/go-redis/redis"
	"github.com/samuel/go-zookeeper/zk"
	"secondkill/sk-core/service/svc_limit"
	"sync"
)

var (
	Redis       RedisConf
	SecKill     SecKillConf
	MysqlConfig MysqlConf
	TraceConfig TraceConf
	Zk          ZookeeperConf
	Etcd        EtcdConf
)

//redis配置
type RedisConf struct {
	//client是Redis客户端，代表零个或多个基础连接池。 对于多个goroutine并发使用是安全的。
	RedisConn            *redis.Client //链接
	Host                 string
	Password             string
	Db                   int
	Proxy2layerQueueName string //队列名称
	Layer2ProxyQueueName string //队列名称
	IdBlackListHash      string //用户黑名单hash表
	IpBlackListHash      string //IP黑名单hash表
	IdBlackListQueue     string //用户黑名单队列
	IpBlackListQueue     string //IP黑名单队列

}

//Etcd配置

type EtcdConf struct {
	EtcdConn          *clientv3.Client //链接
	EtcdSecProductKey string           //商品键
	Host              string
}

//秒杀商品配置信息
type SecKillConf struct {
	RedisConf *RedisConf //redis配置

	CookieSecretKey string   //cookie秘钥
	ReferWhiteList  []string //包名单

	AccessLimitConf AccessLimitConf

	RWBlackLock                  sync.RWMutex
	WriteProxy2LayerGoroutineNum int
	ReadProxy2LayerGoroutineNum  int

	IPBlackMap map[string]bool
	IDBlackMap map[int]bool //用户黑名单

	SecProductInfoMap map[int]*SecProductInfoConf //缓存中的商品信息

	AppWriteToHandleGoroutineNum  int
	AppReadFromHandleGoroutineNum int

	CoreReadRedisGoroutineNum  int
	CoreWriteRedisGoroutineNum int
	CoreHandleGoroutineNum     int

	AppWaitResultTimeOut int

	CoreWaitResultTimeOut int

	MaxRequestWaitTimeout int //最大请求超时时间

	SendToWriteChanTimeout int
	SendToHandChanTimeout  int
	TokenPassWd            string //token秘钥
	//	EtcdConf  *EtcdConf  //EtcdP配置
}

//MySQL配置
type MysqlConf struct {
	Host string
	Port string
	User string
	Pwd  string
	Db   string
}

//Trace配置
type TraceConf struct {
	Host string
	Port string
	Url  string
}

//zookeeper配置
type ZookeeperConf struct {
	ZkConn        *zk.Conn
	SecProductKey string //商品键
}

//商品信息配置
type SecProductInfoConf struct {
	ProductID         int                 `json:"product_id"`           //商品ID
	StartTime         int64               `json:"start_time"`           //开始时间
	EndTime           int64               `json:"end_time"`             //结束时间
	Status            int                 `json:"status"`               //状态
	Total             int                 `json:"total"`                //商品总数
	Left              int                 `json:"left"`                 //商品剩余数量
	OnePersonBuyLimit int                 `json:"one_person_buy_limit"` //单个用户购买数量限制
	BuyRate           float64             `json:"buy_rate"`             //购买频率限制
	SoldMaxlimit      int                 `json:"sold_maxlimit"`        //销售最大量
	SecLimit          *svc_limit.SecLimit `json:"sec_limit"`            //限速控制
}

//访问限制
type AccessLimitConf struct {
	IPSecAccessLimit   int //IP每秒访问限制
	UserSecAccessLimit int //用户每秒范文限制
	IPMinAccessLimit   int //IP每分钟访问限制
	UserMinAccessLimit int //用户每分钟访问限制
}
