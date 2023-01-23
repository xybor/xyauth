package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/xybor/xyauth/internal/logger"
	"github.com/xybor/xyauth/internal/token"
)

type BearerTokenParam struct {
	AccessToken string `json:"access_token" binding:"required"`
}

// VerifyAccessToken verifies the access token in cookies, then adds it into the
// context.
func VerifyAccessToken(ctx *gin.Context) {
	params := BearerTokenParam{}
	if err := ctx.ShouldBindBodyWith(&params, binding.JSON); err != nil {
		if params.AccessToken, err = ctx.Cookie("access_token"); err != nil {
			return
		}
	}

	accessToken := token.AccessToken{}
	if err := token.Verify(params.AccessToken, &accessToken); err != nil {
		logger.Event("access-token-invalid").
			Field("token", params.AccessToken).Field("error", err).Debug()
		return
	}

	ctx.Set("accessToken", accessToken)
}
