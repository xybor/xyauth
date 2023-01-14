package service

import (
	"context"
	"errors"
	"time"

	"github.com/xybor/xyauth/internal/database"
	"github.com/xybor/xyauth/internal/models"
	"github.com/xybor/xyauth/internal/utils"
	"github.com/xybor/xyauth/pkg/token"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CheckWhitelistRefreshToken(t string) error {
	err := database.GetMongoCollection(models.RefreshToken{}).FindOne(
		context.Background(),
		bson.D{{Key: "token", Value: t}},
	).Err()

	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			utils.GetLogger().Event("check-whitelist-refresh-token-failed").
				Field("error", err).Warning()
		}
		return NotFoundError.New("refresh token is revoked or never issued")
	}
	return nil
}

func RevokeRefreshToken(t string) error {
	_, err := database.GetMongoCollection(models.RefreshToken{}).DeleteOne(
		context.Background(),
		bson.D{{Key: "token", Value: t}},
	)

	if err != nil {
		utils.GetLogger().Event("revoke-refresh-token-failed").
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
		utils.GetLogger().Event("janitor-refresh-token-failed").Field("error", err).Warning()
	} else {
		utils.GetLogger().Event("janitor-refresh-token").Field("deleted_count", r.DeletedCount).Info()
	}

	time.AfterFunc(token.RefreshTokenExpiration, JanitorRefreshToken)
}
