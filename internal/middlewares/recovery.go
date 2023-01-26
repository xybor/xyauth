package middlewares

import "github.com/gin-gonic/gin"

func Recovery(ctx *gin.Context) {
	gin.Recovery()(ctx)
}
