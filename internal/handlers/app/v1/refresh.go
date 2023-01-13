package v1

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xybor/xyauth/internal/utils"
	"github.com/xybor/xyauth/pkg/service"
	"github.com/xybor/xyauth/pkg/token"
)

func redirectToLogin(ctx *gin.Context) {
	ctx.Redirect(http.StatusTemporaryRedirect, "/login?redirect_uri="+ctx.Query("redirect_uri"))
}

func RefreshHandler(ctx *gin.Context) {
	cookie, err := ctx.Cookie("refresh_token")
	if err != nil {
		redirectToLogin(ctx)
		return
	}

	refreshToken := token.RefreshToken{}
	if err := token.Verify(cookie, &refreshToken); err != nil {
		utils.GetLogger().Event("refresh-token-invalid").
			Field("cookie", cookie).Field("error", err).Debug()
		redirectToLogin(ctx)
		return
	}

	if err := service.CheckWhitelist(cookie); err != nil {
		if !errors.Is(err, service.NotFoundError) {
			utils.GetLogger().Event("check-whitelist-refresh-token-failed").
				Field("cookie", cookie).Field("error", err).Debug()
		}
		redirectToLogin(ctx)
		return
	}

	value, err := service.CreateAccessToken(refreshToken.Email)
	if err != nil {
		utils.GetLogger().Event("create-access-token-failed").
			Field("email", refreshToken.Email).Field("error", err).Warning()
		redirectToLogin(ctx)
		return
	}

	ctx.SetCookie(
		"access_token", value,
		int(token.AccessTokenExpiration/time.Second), "/",
		utils.GetConfig().MustGet("server.domain").MustString(),
		true, true,
	)

	uri, ok := ctx.GetQuery("redirect_uri")
	if !ok {
		uri = ""
	}
	ctx.Redirect(http.StatusMovedPermanently, uri)
}
