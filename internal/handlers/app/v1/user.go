package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xybor-x/xyerror"
	"github.com/xybor-x/xypriv"
	"github.com/xybor/xyauth/internal/logger"
	"github.com/xybor/xyauth/internal/models"
	"github.com/xybor/xyauth/internal/utils"
	"github.com/xybor/xyauth/pkg/service"
)

func ProfileHandler(ctx *gin.Context) {
	if accessToken, ok := utils.GetAccessToken(ctx); ok {
		ctx.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/user/%s", accessToken.Username))
		return
	}
	utils.RedirectToRefresh(ctx)
}

func UserGETHandler(ctx *gin.Context) {
	var cols []string
	username := ctx.Param("username")
	resource := models.User{Username: username}
	checker := utils.Check(ctx)
	for col := range models.ReadableUserCols {
		if checker.Perform("read", col).On(resource) == nil {
			cols = append(cols, col)
		}
	}

	if checker.Perform("read").On(xypriv.AbstractResource("user_role")) == nil {
		cols = append(cols, "role")
	}

	if len(cols) == 0 {
		utils.HTMLOrLogin(ctx, http.StatusNotFound, "error.html", gin.H{
			"message": "404 Not Found",
		})
		return
	}

	user, err := service.GetUser(username, cols...)
	if err != nil {
		ctx.HTML(http.StatusNotFound, "error.html", gin.H{
			"message": fmt.Sprintf("404 Not Found (%s)", xyerror.Message(err)),
		})
		return
	}

	ctx.HTML(http.StatusOK, "profile.html", gin.H{
		"id":           user.ID,
		"username":     user.Username,
		"email":        user.Email,
		"first_name":   user.FirstName,
		"last_name":    user.LastName,
		"address":      user.Address,
		"phone_number": user.PhoneNumber,
		"role":         user.Role,
	})
}

type UserUpdateParam struct {
	Username    string `form:"username"`
	FirstName   string `form:"first_name"`
	LastName    string `form:"last_name"`
	Address     string `form:"address"`
	PhoneNumber string `form:"phone_number"`
	Role        string `form:"role"`
}

func (p UserUpdateParam) IsEmpty(key string) bool {
	switch key {
	case "username":
		return p.Username == ""
	case "first_name":
		return p.FirstName == ""
	case "last_name":
		return p.LastName == ""
	case "address":
		return p.Address == ""
	case "phone_number":
		return p.PhoneNumber == ""
	case "role":
		return p.Role == ""
	}

	logger.Event("invalid-parameter").Field("key", key).Panic()
	return false
}

func (p UserUpdateParam) Convert() models.User {
	return models.User{
		Username:    p.Username,
		FirstName:   p.FirstName,
		LastName:    p.LastName,
		Address:     p.Address,
		PhoneNumber: p.PhoneNumber,
		Role:        p.Role,
	}
}

func UserPOSTHandler(ctx *gin.Context) {
	param := new(UserUpdateParam)
	if err := ctx.ShouldBind(param); err != nil {
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{"message": "400 Bad Request"})
		return
	}

	checker := utils.Check(ctx)
	username := ctx.Param("username")
	resource := models.User{Username: username}
	for _, col := range models.EditableUserCols {
		if !param.IsEmpty(col) {
			if col != "role" && checker.Perform("update", col).On(resource) == nil {
				continue
			}
			if col == "role" && checker.Perform("update").On(xypriv.AbstractResource("user_role")) == nil {
				continue
			}
			utils.HTMLOrLogin(ctx, http.StatusForbidden, "error.html", gin.H{
				"message": fmt.Sprintf("403 Forbidden (you can not update %s)", col)})
			return
		}
	}

	if err := service.UpdateUser(username, param.Convert()); err != nil {
		logger.Event("update-user-failed", ctx).Field("error", err).Field("username", username).Warning()
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"message": "500 Internal Server Error"})
		return
	}

	ctx.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/user/%s", username))
}
