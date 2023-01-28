package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xybor/xyauth/internal/logger"
	"github.com/xybor/xyauth/internal/token"
)

type BearerTokenParam struct {
	AccessToken string `header:"Authorization" binding:"required"`
}

// VerifyAccessToken verifies the access token in cookies, then adds it into the
// context.
func VerifyAccessToken(ctx *gin.Context) {
	params := new(BearerTokenParam)
	if err := ctx.ShouldBindHeader(params); err == nil {
		authType, value, found := strings.Cut(params.AccessToken, " ")
		if !found || strings.ToLower(authType) != "bearer" {
			return
		}
		params.AccessToken = value
	} else {
		if params.AccessToken, err = ctx.Cookie("access_token"); err != nil {
			return
		}
	}

	accessToken := token.AccessToken{}
	if err := token.Verify(params.AccessToken, &accessToken); err != nil {
		logger.Event("access-token-invalid", ctx).
			Field("token", params.AccessToken).Field("error", err).Debug()
		return
	}

	logger.Event("validated-access-token", ctx).
		Field("request_id", ctx.MustGet("requestID")).
		Field("id", accessToken.ID).Info()
	ctx.Set("accessToken", accessToken)
}
