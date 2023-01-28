package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xybor/xyauth/internal/logger"
	"github.com/xybor/xyauth/internal/utils"
)

func Logger(ctx *gin.Context) {
	requestID := utils.GenerateRandomString(10)
	ctx.Set("requestID", requestID)

	startTime := time.Now()
	ctx.Next()

	logger.Event("http-handle-info", ctx).
		Field("time_taken", time.Since(startTime)).
		Field("code", ctx.Writer.Status()).
		Field("source_ip", ctx.ClientIP()).
		Field("method", ctx.Request.Method).
		Field("path", ctx.Request.URL.Path).
		Info()
}
