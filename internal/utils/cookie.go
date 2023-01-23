package utils

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xybor/xyauth/internal/config"
)

func SetCookie(ctx *gin.Context, name, value string, expiration time.Duration) {
	exp := int(expiration / time.Second)
	if expiration == -1 {
		exp = -1
	}

	ctx.SetCookie(
		name, value,
		exp, "/",
		config.MustGet("server.domain").MustString(),
		true, true,
	)
}
