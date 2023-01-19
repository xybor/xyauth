package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xybor/xyauth/internal/config"
	"github.com/xybor/xyauth/internal/logger"
	"github.com/xybor/xyauth/internal/utils"
	"github.com/xybor/xyauth/pkg/service"
	"github.com/xybor/xyauth/pkg/token"
)

func LogoutHandler(ctx *gin.Context) {
	defer func() {
		ctx.SetCookie(
			"access_token", "", -1, "/",
			config.MustGet("server.domain").MustString(),
			true, true)

		ctx.SetCookie(
			"refresh_token", "", -1, "/",
			config.MustGet("server.domain").MustString(),
			true, true)

		ctx.Redirect(http.StatusTemporaryRedirect, "/")
	}()

	refreshTokenCookie, err := ctx.Cookie("refresh_token")
	if err != nil {
		return
	}

	refreshToken := token.RefreshToken{}
	if err := token.Verify(refreshTokenCookie, &refreshToken); err != nil {
		return
	}

	if utils.Check(ctx).Perform("delete").On(refreshToken) != nil {
		return
	}

	if err := service.RevokeRefreshToken(refreshTokenCookie); err != nil {
		if !errors.Is(err, service.NotFoundError) {
			logger.Event("revoke-refresh-token-failed").Field("error", err).Warning()
		}
		return
	}
}
