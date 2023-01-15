package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xybor/xyauth/internal/utils"
	"github.com/xybor/xyauth/pkg/service"
	"github.com/xybor/xyauth/pkg/token"
)

func LogoutHandler(ctx *gin.Context) {
	defer func() {
		ctx.SetCookie(
			"access_token", "", -1, "/",
			utils.GetConfig().MustGet("server.domain").MustString(),
			true, true)

		ctx.SetCookie(
			"refresh_token", "", -1, "/",
			utils.GetConfig().MustGet("server.domain").MustString(),
			true, true)

		ctx.Redirect(http.StatusTemporaryRedirect, "/")
	}()

	accessTokenCookie, err := ctx.Cookie("access_token")
	if err != nil {
		return
	}

	accessToken := token.AccessToken{}
	if err := token.Verify(accessTokenCookie, &accessToken); err != nil {
		return
	}

	refreshTokenCookie, err := ctx.Cookie("refresh_token")
	if err != nil {
		return
	}

	refreshToken := token.RefreshToken{}
	if err := token.Verify(refreshTokenCookie, &refreshToken); err != nil {
		return
	}

	if accessToken.Email != refreshToken.Email {
		return
	}

	if err := service.RevokeRefreshToken(refreshTokenCookie); err != nil {
		if !errors.Is(err, service.NotFoundError) {
			utils.GetLogger().Event("revoke-refresh-token-failed").Field("error", err).Warning()
		}
		return
	}
}
