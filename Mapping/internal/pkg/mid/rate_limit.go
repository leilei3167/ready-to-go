package mid

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"net/http"
	"time"
)

// RateLimitMiddleware 实现一个限流中间件
func RateLimitMiddleware(fillInterval time.Duration, cap, quantum int64) gin.HandlerFunc {
	bucket := ratelimit.NewBucketWithQuantum(fillInterval, cap, quantum)
	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) < 1 { //没有可用的令牌时,直接告知服务器忙
			c.String(http.StatusForbidden, "服务器繁忙...")
			c.Abort()
			return
		}
		c.Next()
	}
}

//后续可增加其他中间件...
