package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xybor-x/xyerror"
	"github.com/xybor/xyauth/pkg/service"
)

type Parameters struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

func RegisterGETHandler(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "register.html", nil)
}

func RegisterPOSTHandler(ctx *gin.Context) {
	params := new(Parameters)
	if err := ctx.ShouldBind(params); err != nil {
		ctx.HTML(http.StatusBadRequest, "register.html", gin.H{
			"message": "Invalid request",
		})
		return
	}

	err := service.Register(params.Email, params.Password, "")
	switch {
	case errors.Is(err, service.FormatError):
		ctx.HTML(http.StatusBadRequest, "register.html", gin.H{
			"message": xyerror.Message(err),
		})
	case errors.Is(err, service.DuplicatedError):
		ctx.HTML(http.StatusConflict, "register.html", gin.H{
			"message": xyerror.Message(err),
		})
	case err != nil:
		ctx.HTML(http.StatusInternalServerError, "register.html", gin.H{
			"message": "Something is wrong, please contact to administrator if the issue persists",
		})
	default:
		ctx.HTML(http.StatusOK, "register.html", gin.H{
			"message": "You registered successfully",
		})
	}
}
