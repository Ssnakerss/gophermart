package db

import (
	"context"
	"errors"
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

type DB struct {
	ConString string
	GormDB    *gorm.DB
}

func New(conString string, logLevel LogLevel) *DB {
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

	var newGorm DB
	newGorm.ConString = conString
	var err error
	newGorm.GormDB, err = gorm.Open(postgres.Open(newGorm.ConString), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatal(err)
	}
	return &newGorm
}

func (db *DB) Migrate(ctx context.Context) error {
	return db.GormDB.WithContext(ctx).AutoMigrate(&models.Order{})
}

func (db *DB) SaveOrder(ctx context.Context, order *models.Order) error {
	return db.GormDB.WithContext(ctx).Save(order).Error
}

func (db *DB) GetOrder(ctx context.Context, order *models.Order) error {
	err := db.GormDB.WithContext(ctx).First(order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.ErrRecordNotFound
	}
	return err
}

func (db *DB) GetAllOrders(ctx context.Context, order *models.Order) []models.Order {
	var orders []models.Order
	db.GormDB.WithContext(ctx).Order("time_stamp DESC").Find(&orders, order)
	return orders
}
