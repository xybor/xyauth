package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xybor-x/xyerror"
	"github.com/xybor/xyauth/internal/logger"
	"github.com/xybor/xyauth/internal/token"
	"github.com/xybor/xyauth/internal/utils"
	"github.com/xybor/xyauth/pkg/service"
)

func redirectToLogin(ctx *gin.Context) {
	ctx.Redirect(http.StatusMovedPermanently, "/login?redirect_uri="+ctx.Query("redirect_uri"))
}

func RefreshHandler(ctx *gin.Context) {
	cookie, err := ctx.Cookie("refresh_token")
	if err != nil {
		redirectToLogin(ctx)
		return
	}

	refreshToken := token.RefreshToken{}
	if err := token.Verify(cookie, &refreshToken); err != nil {
		logger.Event("refresh-token-invalid", ctx).
			Field("cookie", cookie).Field("error", err).Debug()
		redirectToLogin(ctx)
		return
	}

	newRefreshToken, err := service.InheritRefreshToken(refreshToken)
	if err != nil {
		if !errors.Is(err, service.SecurityError) {
			logger.Event("inherit-refresh-token-failed", ctx).
				Field("user_id", refreshToken.ID).Field("error", err).Warning()
		}

		utils.SetCookie(ctx, "access_token", "", -1)
		utils.SetCookie(ctx, "refresh_token", "", -1)
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"message": fmt.Sprintf("500 Internal Server Error (%s)", xyerror.Message(err)),
		})
		ctx.Abort()
		return
	}

	accessToken, err := service.CreateAccessToken(refreshToken.ID)
	if err != nil {
		logger.Event("create-access-token-failed", ctx).
			Field("user_id", refreshToken.ID).Field("error", err).Warning()
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"message": "500 Internal Server Error (can not create access token)",
		})
		return
	}

	utils.SetCookie(ctx, "access_token", accessToken, token.AccessTokenExpiration)
	utils.SetCookie(ctx, "refresh_token", newRefreshToken, token.RefreshTokenExpiration)

	uri, ok := ctx.GetQuery("redirect_uri")
	if !ok {
		uri = "/"
	}
	ctx.Redirect(http.StatusTemporaryRedirect, uri)
}
