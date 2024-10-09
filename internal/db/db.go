package db

import (
	"context"
	"errors"
	"log"

	"github.com/Ssnakerss/gophermart/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const ConString = "host=localhost user=postgres dbname=postgres password=postgres sslmode=disable"

type Db struct {
	conString string
	db        *gorm.DB
}

func New(conString string) *Db {
	//`postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable`
	//"host=localhost user=postgres dbname=postgres password=postgres sslmode=disable"
	var db Db
	db.conString = conString
	var err error
	db.db, err = gorm.Open(postgres.Open(db.conString), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return &db
}

func (db *Db) Migrate(ctx context.Context) error {
	return db.db.WithContext(ctx).AutoMigrate(&models.Order{})
}

func (db *Db) SaveOrder(ctx context.Context, order *models.Order) error {
	return db.db.WithContext(ctx).Save(order).Error
}

func (db *Db) GetOrder(ctx context.Context, order *models.Order) error {
	err := db.db.WithContext(ctx).First(order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.ErrRecordNotFound
	}
	return err
}

func (db *Db) GetAllOrders(ctx context.Context, order *models.Order) []models.Order {
	var orders []models.Order
	db.db.WithContext(ctx).Find(&orders)
	return orders
}
