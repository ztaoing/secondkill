/**
* @Author:zhoutao
* @Date:2020/7/5 上午10:11
 */

package config

import (
	"github.com/go-kit/kit/log"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/spf13/viper"
	"os"
	"secondkill/pkg/bootstrap"
	"secondkill/pkg/config"
	"secondkill/sk-app/model"
	"sync"
)

const KConfigType = "CONFIG_TYPE"
const (
	ProductStatusNormal       = 0 //商品状态正常
	ProductStatusSaleOut      = 1 //商品卖光
	ProductStatusForceSaleOut = 2 //商品强制卖光
)

var ZipkinTracer *zipkin.Tracer
var Logger log.Logger

func init() {

	Logger = log.NewLogfmtLogger(os.Stderr)
	Logger = log.With(Logger, "timestamp", log.DefaultTimestamp)
	Logger = log.With(Logger, "caller", log.DefaultCaller)

	viper.AutomaticEnv()
	initDefault()

	if err := config.LoadRemoteConfig(); err != nil {
		Logger.Log("faild to load remote config", err)
	}
	if err := config.Sub("mysql", &config.MysqlConfig); err != nil {
		Logger.Log("faild to parse mysql", err)
	}
	if err := config.Sub("service", &config.SecKill); err != nil {
		Logger.Log("faild to parse service seckill", err)
	}
	if err := config.Sub("redis", &config.Redis); err != nil {
		Logger.Log("faild to parse reids", err)
	}

	zipkinURL := "http://" + config.TraceConfig.Host + ":" + config.TraceConfig.Port + config.TraceConfig.Url
	Logger.Log("zipkin url", zipkinURL)
	initTracer(zipkinURL)
}

func initDefault() {
	viper.SetDefault(KConfigType, "yaml")
}

func initTracer(zipkinURL string) {
	var (
		err           error
		useNoopTracer = (zipkinURL == "")
		reporter      = zipkinhttp.NewReporter(zipkinURL)
	)
	zEndpoint, _ := zipkin.NewEndpoint(bootstrap.DiscoverConfig.ServiceName, bootstrap.HttpConfig.Port)
	_, err = zipkin.NewTracer(reporter,
		zipkin.WithLocalEndpoint(zEndpoint),
		zipkin.WithNoopTracer(useNoopTracer),
	)
	if err != nil {
		Logger.Log("err", err)
		os.Exit(1)
	}
	if !useNoopTracer {
		Logger.Log("tracer", "zipkin", "type", "nativa", "url", zipkinURL)
	}

}

/**
秒杀
*/

type SkAppCtx struct {
	//请求
	SecReqChan       chan *model.SecRequest
	SecReqChanSize   int
	RWSecProductLock sync.RWMutex

	UserConnMap     map[string]chan *model.SecResult
	UserConnMapLock sync.Mutex
}

var SkAppContext = &SkAppCtx{
	//用户连接
	UserConnMap: make(map[string]chan *model.SecResult, 1024),
	//请求
	SecReqChan: make(chan *model.SecRequest, 1024),
}
