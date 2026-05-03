package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	limiter "github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"

	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
)

func RateLimiter(rateStr string) gin.HandlerFunc {
	rate, err := limiter.NewRateFromFormatted(rateStr)
	store := memory.NewStore()
	instance := limiter.New(store, rate)

	if err != nil {
		fmt.Print("hello from the otherside")
	}

	return mgin.NewMiddleware(instance)
}
