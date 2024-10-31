package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Ssnakerss/gophermart/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const ConString = "host=localhost user=postgres dbname=postgres password=postgres sslmode=disable"

type LogLevel int

const (
	// Silent silent log level
	Silent LogLevel = iota + 1
	// Error error log level
	Error
	// Warn warn log level
	Warn
	// Info info log level
	Info
)

type GormDB struct {
	ConString string
	DB        *gorm.DB
}

func New(conString string, logLevel LogLevel) (*GormDB, error) {
	//`postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable`
	//"host=localhost user=postgres dbname=postgres password=postgres sslmode=disable"
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,               // Slow SQL threshold
			LogLevel:                  logger.LogLevel(logLevel), // Log level
			IgnoreRecordNotFoundError: true,                      // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,                      // Disable color
		},
	)
	var newGorm GormDB

	newGorm.ConString = conString
	var err error
	newGorm.DB, err = gorm.Open(postgres.Open(newGorm.ConString), &gorm.Config{
		Logger: newLogger,
	})

	return &newGorm, err
}

func (db *GormDB) Migrate(ctx context.Context) error {
	return db.DB.WithContext(ctx).AutoMigrate(
		&models.Order{},
		&models.Account{},
		&models.Transaction{},
		&models.User{},
	)
}
