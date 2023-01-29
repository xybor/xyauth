package database

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"strings"

	"github.com/xybor-x/xyerror"
	"github.com/xybor/xyauth/internal/certificate"
	"github.com/xybor/xyauth/internal/config"
	"github.com/xybor/xyauth/internal/logger"
	"github.com/xybor/xyauth/internal/models"
	"github.com/xybor/xyauth/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoDB *mongo.Database
var mongoCollections = make(map[string]*mongo.Collection)

func InitMongoDB() error {
	scheme := config.GetDefault("mongodb.scheme", "mongodb").MustString()
	host := config.GetDefault("mongodb.host", "localhost").MustString()
	port := config.GetDefault("mongodb.port", 27017).MustInt()

	dbName := config.MustGet("MONGO_INITDB_DATABASE").MustString()
	username := config.MustGet("MONGO_INITDB_ROOT_USERNAME").MustString()
	password := config.MustGet("MONGO_INITDB_ROOT_PASSWORD").MustString()

	dsn := fmt.Sprintf("%s://%s:%s@%s:%d",
		scheme,
		username, password,
		host, port,
	)

	parameters := make([]string, 0)
	for _, key := range []string{"tls", "replicaSet", "readPreference", "retryWrites"} {
		if val, ok := config.Get("mongodb." + key); ok {
			parameters = append(parameters, key+"="+val.MustString())
		}
	}

	if p := strings.Join(parameters, "&"); p != "" {
		dsn += "/?" + p
	}

	opts := options.Client().ApplyURI(dsn)

	if val, ok := config.Get("MONGO_SSL_CA_CERTS"); ok {
		tlsConfig, err := getCustomTLSConfig(val.MustString())
		if err != nil {
			return err
		}
		opts.SetTLSConfig(tlsConfig)
	}

	client, err := mongo.NewClient(opts)
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
	logger.Event("connect-to-mongodb").
		Field("host", host).Field("port", port).Field("db", dbName).Info()

	return nil
}

func GetMongoCollection(a any) *mongo.Collection {
	c, err := utils.GetSnakeCase(a)
	if err != nil || c == "" {
		logger.Event("extract-snake-case-failed").
			Field("error", err).Field("struct", a).Panic()
	}

	if _, ok := mongoCollections[c]; !ok {
		mongoCollections[c] = mongoDB.Collection(c)
	}

	return mongoCollections[c]
}

func DropAllNoSQL() error {
	for _, a := range models.AllNoSQL {
		if err := GetMongoCollection(a).Drop(context.Background()); err != nil {
			return err
		}
	}
	return nil
}

func getCustomTLSConfig(caFile string) (*tls.Config, error) {
	certs, err := certificate.GetCertificateContent(caFile)
	if err != nil {
		return nil, err
	}

	tlsConfig := new(tls.Config)
	tlsConfig.RootCAs = x509.NewCertPool()
	if ok := tlsConfig.RootCAs.AppendCertsFromPEM(certs); !ok {
		return tlsConfig, xyerror.ValueError.New("Failed parsing pem file")
	}

	return tlsConfig, nil
}
