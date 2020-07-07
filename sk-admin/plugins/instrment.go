/**
* @Author:zhoutao
* @Date:2020/7/7 上午9:54
 */

package plugins

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/juju/ratelimit"
	"golang.org/x/time/rate"
)

var ErrLimitExceed = errors.New("rate limit exceed")

//juju/ratelimite 创建限流中间件
func NewTokenBucketLimitterWithJuju(bkt *ratelimit.Bucket) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			//限流控制
			if bkt.TakeAvailable(1) == 0 {
				return nil, ErrLimitExceed
			}
			//经过限流控制
			return next(ctx, request)
		}
	}
}

//time/rate 穿件限流组件
func NewTokenBucketLimitterWithBuidIn(bkt *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !bkt.Allow() {
				return nil, ErrLimitExceed
			}
			//经过限流控制
			return next(ctx, request)
		}
	}
}
