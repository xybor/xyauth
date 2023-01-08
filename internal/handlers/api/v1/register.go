package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xybor-x/xyerror"
	"github.com/xybor/xyauth/pkg/service"
)

type RegisterParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func RegisterHandler(ctx *gin.Context) {
	params := new(RegisterParams)
	ctx.ShouldBind(params)

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
