package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xybor-x/xyerror"
	"github.com/xybor/xyauth/internal/logger"
	"github.com/xybor/xyauth/internal/token"
	"github.com/xybor/xyauth/pkg/service"
)

type RefreshParams struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func RefreshHandler(ctx *gin.Context) {
	params := new(RefreshParams)
	if err := ctx.ShouldBind(params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid parameters"})
		return
	}

	refreshToken := token.RefreshToken{}
	if err := token.Verify(params.RefreshToken, &refreshToken); err != nil {
		if errors.Is(err, token.TokenError) {
			ctx.JSON(http.StatusForbidden, gin.H{"message": xyerror.Message(err)})
		} else {
			logger.Event("refresh-token-invalid").
				Field("token", params.RefreshToken).Field("error", err).Warning()
			ctx.Status(http.StatusInternalServerError)
		}
		return
	}

	newRefreshToken, err := service.InheritRefreshToken(refreshToken)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"message": xyerror.Message(err)})
		return
	}

	accessToken, err := service.CreateAccessToken(refreshToken.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": xyerror.Message(err)})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
	})
}
