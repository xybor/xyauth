package service

import (
	"context"
	"time"

	"github.com/xybor/xyauth/internal/database"
	"github.com/xybor/xyauth/internal/logger"
	"github.com/xybor/xyauth/internal/models"
	"github.com/xybor/xyauth/internal/token"
	"github.com/xybor/xyauth/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateRefreshToken(id uint) (string, error) {
	family := utils.GenerateRandomString(32)

	value, err := token.Create(token.NewRefreshTokenConfig(token.RefreshToken{
		ID: id, Family: family, FamilyID: 0,
	}))
	if err != nil {
		return "", err
	}

	_, err = database.GetMongoCollection(models.RefreshToken{}).InsertOne(
		context.Background(), models.RefreshToken{
			ID:         id,
			Family:     family,
			Counter:    0,
			Expiration: time.Now().Add(token.RefreshTokenExpiration),
		},
	)

	if err != nil {
		logger.Event("whitelist-refresh-token-failed").
			Field("token", value).Field("error", err).Error()
		return "", ServiceError.New("can not insert refresh token")
	}

	return value, nil
}

func InheritRefreshToken(request token.RefreshToken) (string, error) {
	result := database.GetMongoCollection(models.RefreshToken{}).FindOne(
		context.Background(), bson.D{
			{Key: "family", Value: request.Family},
			{Key: "id", Value: request.ID},
		},
	)

	if result.Err() != nil {
		return "", NotFoundError.New("refresh token is expired or invalid")
	}

	info := models.RefreshToken{}
	if err := result.Decode(&info); err != nil {
		logger.Event("whitelist-refresh-token-failed").
			Field("result", result).Field("error", err).Error()
		return "", ValueError.New("invalid whitelist token")
	}

	if request.FamilyID < info.Counter {
		RevokeRefreshToken(request)
		return "", SecurityError.New("the request token might be stolen")
	}

	updateResult, err := database.GetMongoCollection(models.RefreshToken{}).UpdateOne(
		context.Background(), bson.D{
			{Key: "family", Value: request.Family},
			{Key: "id", Value: request.ID},
			{Key: "counter", Value: info.Counter}, // Avoid race condition
		},
		bson.D{{Key: "$inc", Value: bson.D{{Key: "counter", Value: 1}}}},
	)

	if err != nil {
		logger.Event("update-refresh-token-failed").
			Field("family", info.Family).Field("error", err).Error()
		return "", ServiceError.New("can not update refresh token")
	}

	if updateResult.MatchedCount == 0 {
		return "", SecurityError.New("may be a race condition occurred")
	}

	value, err := token.Create(token.NewRefreshTokenConfig(token.RefreshToken{
		ID:       info.ID,
		Family:   info.Family,
		FamilyID: info.Counter + 1,
	}))
	if err != nil {
		return "", err
	}

	return value, nil
}

func RevokeRefreshToken(request token.RefreshToken) error {
	_, err := database.GetMongoCollection(models.RefreshToken{}).DeleteOne(
		context.Background(),
		bson.D{
			{Key: "family", Value: request.Family},
			{Key: "id", Value: request.ID},
		},
	)

	if err != nil {
		logger.Event("revoke-refresh-token-failed").
			Field("error", err).Warning()
		return NotFoundError.New("refresh token can not be revoked")
	}
	return nil
}

func JanitorRefreshToken() {
	r, err := database.GetMongoCollection(models.RefreshToken{}).DeleteMany(
		context.Background(),
		bson.D{{Key: "expiration", Value: bson.D{{Key: "$lt", Value: time.Now()}}}},
	)

	if err != nil {
		logger.Event("janitor-refresh-token-failed").Field("error", err).Warning()
	} else {
		logger.Event("janitor-refresh-token").Field("deleted_count", r.DeletedCount).Info()
	}

	time.AfterFunc(token.RefreshTokenExpiration, JanitorRefreshToken)
}
