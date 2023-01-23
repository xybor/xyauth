package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xybor-x/xyerror"
	"github.com/xybor-x/xypriv"
	"github.com/xybor/xyauth/internal/utils"
	"github.com/xybor/xyauth/pkg/service"
)

type RegisterParams struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role"`
}

func RegisterHandler(ctx *gin.Context) {
	params := new(RegisterParams)
	if err := ctx.ShouldBind(params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid parameters"})
		return
	}

	table := xypriv.AbstractResource("user_credential_table")
	if utils.Check(ctx).Perform("create", params.Role).On(table) != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "permission denied"})
		return
	}

	err := service.Register(params.Email, params.Password, params.Role)
	switch {
	case errors.Is(err, service.DuplicatedError):
		ctx.JSON(http.StatusConflict, gin.H{"message": xyerror.Message(err)})
	case errors.Is(err, service.FormatError):
		ctx.JSON(http.StatusBadRequest, gin.H{"message": xyerror.Message(err)})
	case err != nil:
		ctx.Status(http.StatusInternalServerError)
	default:
		ctx.Status(http.StatusOK)
	}
}
