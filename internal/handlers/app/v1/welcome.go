package v1

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xybor-x/xypriv"
	"github.com/xybor/xyauth/internal/utils"
)

func WelcomeHandler(ctx *gin.Context) {
	page := xypriv.AbstractResource("welcome_page")
	if utils.Check(ctx).Perform("read").On(page) != nil {
		utils.HTMLOrLogin(ctx, http.StatusForbidden, "error.html", gin.H{
			"message": "403 Forbidden (your account doesn't have the permission to access website)",
		})
		return
	}

	name := "Stranger"
	if accessToken, ok := utils.GetAccessToken(ctx); ok {
		name = strings.Split(accessToken.GetUser().Email, "@")[0]
	}
	ctx.HTML(http.StatusOK, "welcome.html", gin.H{"name": name})
}
