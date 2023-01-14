package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xybor-x/xyerror"
	"github.com/xybor/xyauth/internal/utils"
	"github.com/xybor/xyauth/pkg/service"
	"github.com/xybor/xyauth/pkg/token"
)

type RefreshParams struct {
	RefreshToken string `json:"refresh_token"`
}

func RefreshHandler(ctx *gin.Context) {
	params := RefreshParams{}
	err := ctx.ShouldBind(&params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "expect a refresh_token"})
		return
	}

	refreshToken := token.RefreshToken{}
	if err := token.Verify(params.RefreshToken, &refreshToken); err != nil {
		if errors.Is(err, token.TokenError) {
			ctx.JSON(http.StatusForbidden, gin.H{"message": xyerror.Message(err)})
		} else {
			utils.GetLogger().Event("refresh-token-invalid").
				Field("token", params.RefreshToken).Field("error", err).Warning()
			ctx.Status(http.StatusInternalServerError)
		}
		return
	}

	if err := service.CheckWhitelistRefreshToken(params.RefreshToken); err != nil {
		if errors.Is(err, service.NotFoundError) {
			ctx.JSON(http.StatusForbidden, gin.H{"message": xyerror.Message(err)})
		} else {
			utils.GetLogger().Event("check-whitelist-refresh-token-failed").
				Field("token", params.RefreshToken).Field("error", err).Warning()
			ctx.Status(http.StatusInternalServerError)
		}
		return
	}

	value, err := service.CreateAccessToken(refreshToken.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": xyerror.Message(err)})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"access_token": value})
}
