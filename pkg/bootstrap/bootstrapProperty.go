/**
* @Author:zhoutao
* @Date:2020/7/1 上午7:49
 */

package bootstrap

var (
	HttpConfig         HttpConf
	DiscoverConfig     DiscoverConf
	ConfigServerConfig ConfigServerConf
	PrpcConfig         PrpcConf
)

//http配置
type HttpConf struct {
	Host string
	Port string
}

//服务发现注册配置
type DiscoverConf struct {
	Host        string
	Port        string
	ServiceName string
	Weight      int
	InstanceID  string
}

//配置中心配置
type ConfigServerConf struct {
	ID      string //服务名
	Profile string //
	Label   string //
}

//rpc配置
type PrpcConf struct {
	Port string
}
