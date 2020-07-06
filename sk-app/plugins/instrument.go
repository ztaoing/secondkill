/**
* @Author:zhoutao
* @Date:2020/7/5 上午10:45
 */

package plugins

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/juju/ratelimit"
	"golang.org/x/time/rate"
)

var (
	ErrLimitExceed = errors.New("rate limit exceed")
)

/**
juju限流组件
*/

func NewTokenBucketLimitterWithJuju(bkt *ratelimit.Bucket) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if bkt.TakeAvailable(1) == 0 {
				return nil, ErrLimitExceed
			}
			return next(ctx, request)
		}
	}
}

/**
time/rate
*/
func NewTokenBuketLimitterWithBuildIn(bkt *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !bkt.Allow() {
				return nil, ErrLimitExceed
			}
			return next(ctx, request)
		}
	}
}
