package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xybor/xyauth/internal/logger"
)

func Logger(ctx *gin.Context) {
	startTime := time.Now()
	ctx.Next()

	logger.Event("http-handle-info").
		Field("time_taken", time.Since(startTime)).
		Field("code", ctx.Writer.Status()).
		Field("source_ip", ctx.ClientIP()).
		Field("method", ctx.Request.Method).
		Field("path", ctx.Request.URL.Path).
		Info()
}
