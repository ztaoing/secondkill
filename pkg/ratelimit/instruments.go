/**
* @Author:zhoutao
* @Date:2020/7/2 上午9:12
限流组件
*/

/**
在秒杀场景中，用于业务应用系统的负载能力有限，为了防止非预期的请求对系统造成过大的压力而拖垮业务应用系统，每个Api接口都有访问频率上限
API接口的流量控制策略有分流、降级、限流等。
被组件是限流策略，虽然降低了服务接口的访问频率和并发量，却换来了服务接口和业务应用系统的高可用
*/
/**
漏桶算法:主要目的是控制数据注入到系统的速率，平滑应对系统的突发流量，为系统提供一个稳定的请求流量。强行限制数据的传输速率。按固定的速率流出请求。不允许突发流量.
令牌桶算法:按固定的速率往桶中添加令牌。允许突发请求。支持一次拿3-4个令牌
使用rate实现限流：go包，是基于令牌桶实现
*/
package ratelimit

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/time/rate"
	"time"
)

var ErrLimitExceed = errors.New("rate limit exceed!")

//在秒杀实例项目中，endpoint.Middleware 为每个endpoint提供限流功能

//使用x/time/rate创建限流中间件
func NewTokenBucktLimitterWithBuildIn(bkt *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			//如果限流器不放行
			if !bkt.Allow() {
				return nil, ErrLimitExceed
			}
			return next(ctx, request)
		}
	}
}

func DynamicLimitter(interval int, burst int) endpoint.Middleware {
	//NewLimiter返回一个新的限制器，该限制器允许事件的发生率达到r(interval)，并允许突发最多b(burst)个令牌
	bucket := rate.NewLimiter(rate.Every(time.Duration(interval)*time.Second), burst)
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !bucket.Allow() {
				return nil, ErrLimitExceed
			}
			return next(ctx, request)
		}
	}
}
