package middleware

import (
	"backend/pkg/response"
	"fmt"
	"net/http"

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

		return func(c *gin.Context) {
			response.Error(c, http.StatusInternalServerError, fmt.Sprintf("invalid rate format: %v", err))
			c.Abort()
		}
	}

	return mgin.NewMiddleware(instance)
}
