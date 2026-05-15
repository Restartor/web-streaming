package middleware

import (
	"backend/pkg/logger"

	"github.com/gin-gonic/gin"
	limiter "github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"

	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
)

func RateLimiter(rateStr string) gin.HandlerFunc {
	rate, err := limiter.NewRateFromFormatted(rateStr)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("rate limiter config error")
	}

	storing := memory.NewStore()
	instance := limiter.New(storing, rate)

	return mgin.NewMiddleware(instance)
}
