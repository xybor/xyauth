package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/xybor/xyauth/internal/config"
	"github.com/xybor/xyauth/internal/logger"
	"github.com/xybor/xyauth/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"
)

var postgresDB *gorm.DB

// InitPostgresDB inits the db with Gorm. The parameter p is using only for testing.
func InitPostgresDB(p gorm.ConnPool) error {
	loglevel := config.GetDefault("postgresql.loglevel", int(gormlog.Warn)).MustInt()
	newLogger := gormlog.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		gormlog.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  gormlog.LogLevel(loglevel),
			IgnoreRecordNotFoundError: true,
		},
	)

	if p != nil {
		var err error
		postgresDB, err = gorm.Open(postgres.New(postgres.Config{
			Conn: p,
		}), &gorm.Config{Logger: newLogger})
		if err != nil {
			return err
		}
		return nil
	}

	dsnFormat := "host=%s user=%s password=%s dbname=%s port=%s sslmode=%s"
	// Prioritize environment variable over config file.
	host := config.MustGet("postgresql.host").MustString()
	port := config.MustGet("postgresql.port").MustString()
	user := config.MustGet("POSTGRES_USER").MustString()
	password := config.MustGet("POSTGRES_PASSWORD").MustString()
	dbname := config.MustGet("POSTGRES_DB").MustString()
	sslmode := config.GetDefault("postgresql.sslmode", "disable").MustString()
	dsn := fmt.Sprintf(dsnFormat, host, user, password, dbname, port, sslmode)

	timezone, ok := config.Get("postgresql.timezone")
	if ok {
		dsn += fmt.Sprintf(" TimeZone=%s", timezone.MustString())
	}

	var err error
	nRetries := config.GetDefault("postgresql.retries", 3).MustInt()
	retryDuration := config.GetDefault("postgresql.retry_duration", time.Second).MustDuration()
	for i := 0; i < nRetries; i++ {
		postgresDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})
		if err == nil {
			break
		}
		logger.Event("connect-to-database-failed").Field("error", err).Error()
		time.Sleep(retryDuration)
	}

	if err != nil {
		return err
	}

	if err := postgresDB.AutoMigrate(models.AllSQL...); err != nil {
		return err
	}

	logger.Event("connect-to-postgresql").
		Field("host", host).Field("port", port).Info()

	return nil
}

func GetPostgresDB() *gorm.DB {
	return postgresDB
}

func DropAllSQL() error {
	return postgresDB.Migrator().DropTable(models.AllSQL...)
}
