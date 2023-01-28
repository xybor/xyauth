package service

import (
	"errors"

	"github.com/xybor/xyauth/internal/database"
	"github.com/xybor/xyauth/internal/logger"
	"github.com/xybor/xyauth/internal/models"
	"gorm.io/gorm"
)

func UsernameExists(username string) bool {
	result := database.GetPostgresDB().Select("1").Where("username=?", username).Take(&models.User{})
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logger.Event("get-user-info-failed").Field("error", result.Error).Warning()
		}
		return false
	}
	return true
}

func GetUser(username string, cols ...string) (models.User, error) {
	user := models.NewMaskedUser()
	for i := range cols {
		if !user.IsReadable(cols[i]) {
			return user, PermissionError.Newf("col %s is not visible", cols[i])
		}
	}

	result := database.GetPostgresDB().Select(cols).Where("username=?", username).Take(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return user, NotFoundError.Newf("could not find the user %s", username)
		}
		logger.Event("get-user-info-failed").Field("error", result.Error).Warning()
		return user, ServiceError.New("failed to get user info")
	}

	return user, nil
}

func UpdateUser(username string, info models.User) error {
	result := database.GetPostgresDB().
		Model(&models.User{}).
		Where("username=?", username).
		Updates(info)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return NotFoundError.Newf("could not find the user %s", username)
		}
		logger.Event("failed-to-update-profile").Field("error", result.Error).Warning()
		return ServiceError.New("failed to update data")
	}

	return nil
}
