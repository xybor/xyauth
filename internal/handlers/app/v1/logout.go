package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xybor/xyauth/internal/logger"
	"github.com/xybor/xyauth/internal/token"
	"github.com/xybor/xyauth/internal/utils"
	"github.com/xybor/xyauth/pkg/service"
)

func LogoutHandler(ctx *gin.Context) {
	defer func() {
		utils.SetCookie(ctx, "access_token", "", -1)
		utils.SetCookie(ctx, "refresh_token", "", -1)
		ctx.Redirect(http.StatusTemporaryRedirect, "/")
	}()

	val, err := ctx.Cookie("refresh_token")
	if err != nil {
		return
	}

	refreshToken := token.RefreshToken{}
	if err := token.Verify(val, &refreshToken); err != nil {
		return
	}

	if utils.Check(ctx).Perform("delete").On(refreshToken) != nil {
		return
	}

	if err := service.RevokeRefreshToken(refreshToken); err != nil {
		if !errors.Is(err, service.NotFoundError) {
			logger.Event("revoke-refresh-token-failed", ctx).Field("error", err).Warning()
		}
		return
	}
}
