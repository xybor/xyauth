package database

import (
	"context"
	"fmt"

	"github.com/xybor/xyauth/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoDB *mongo.Database
var mongoCollections = make(map[string]*mongo.Collection)

func InitMongoDB() error {
	host := utils.GetConfig().GetDefault("mongodb.host", "localhost").MustString()
	port := utils.GetConfig().GetDefault("mongodb.port", 27017).MustInt()

	dbName := utils.GetConfig().MustGet("MONGO_INITDB_DATABASE").MustString()
	username := utils.GetConfig().MustGet("MONGO_INITDB_ROOT_USERNAME").MustString()
	password := utils.GetConfig().MustGet("MONGO_INITDB_ROOT_PASSWORD").MustString()

	dsn := fmt.Sprintf("mongodb://%s:%s@%s:%d",
		username, password,
		host, port,
	)

	client, err := mongo.NewClient(options.Client().ApplyURI(dsn))
	if err != nil {
		return err
	}

	if err := client.Connect(context.Background()); err != nil {
		return err
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		return err
	}

	mongoDB = client.Database(dbName)

	utils.GetLogger().Event("connect-to-mongodb").
		Field("host", host).Field("port", port).Field("db", dbName).Info()

	return nil
}

func GetMongoCollection(a any) *mongo.Collection {
	c, err := utils.GetSnakeCase(a)
	if err != nil || c == "" {
		utils.GetLogger().Event("extract-snake-case-failed").
			Field("error", err).Field("struct", a).Panic()
	}

	if _, ok := mongoCollections[c]; !ok {
		mongoCollections[c] = mongoDB.Collection(c)
	}

	return mongoCollections[c]
}
