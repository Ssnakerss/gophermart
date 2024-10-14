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

type GormDB struct {
	ConString string
	DB        *gorm.DB
}

func New(conString string, logLevel LogLevel) *GormDB {
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
	if err != nil {
		log.Fatal(err)
	}
	return &newGorm
}

func (db *GormDB) Migrate(ctx context.Context) error {
	return db.DB.WithContext(ctx).AutoMigrate(&models.Order{}, &models.Account{}, &models.Transaction{})
}

func (db *GormDB) SaveOrder(ctx context.Context, order *models.Order) error {
	return db.DB.WithContext(ctx).Save(order).Error
}

func (db *GormDB) GetOrder(ctx context.Context, order *models.Order) error {
	err := db.DB.WithContext(ctx).First(order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.ErrRecordNotFound
	}
	return err
}

func (db *GormDB) GetAllOrders(ctx context.Context, order *models.Order) []models.Order {
	var orders []models.Order
	db.DB.WithContext(ctx).Order("time_stamp DESC").Find(&orders, order)
	return orders
}

// Account related operations
// Сначала пробуем обновить баланс -  важно при списании
// делаем все в транзакции для сохранения целостности данных
func (db *GormDB) PostTransaction(ctx context.Context, accTransaction *models.Transaction) error {

	//TODO сделать проверку на существование аккаунта
	if accTransaction.Indicator != "D" && accTransaction.Indicator != "C" {
		return models.ErrInvalidIndicator
	}
	account := models.Account{
		UserID: accTransaction.UserID,
	}

	col := ""
	colSql := ""

	//TODO при ошибке записать транзакцию в журнал	со статусом E
	var terr error
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var err error
		switch accTransaction.Indicator {
		//пополнение счета
		case "D":
			err = tx.Model(&account).Update("balance", gorm.Expr("balance + ?", accTransaction.Bonus)).Error
			if err != nil {
				return err
			}
			col = "debit"
			colSql = "debit + ?"
		//списание со счета
		case "C":
			result := tx.Model(&account).
				//при списании надо проверить достаочно ли средств на счету
				Where("balance >= ?", accTransaction.Bonus).
				Update("balance", gorm.Expr("balance - ?", accTransaction.Bonus))
				//если списание не прощло - ставим ошибку
			if result.RowsAffected == 0 {
				terr = models.ErrInsufficientFunds
			} else {
				err = result.Error
			}

			col = "credit"
			colSql = "credit + ?"
		}

		if terr == models.ErrInsufficientFunds {
			accTransaction.Indicator = "E"
		}

		if err != nil {
			return err
		}
		//если транзакция обновила баланс успешно - обноаляем и пола счета с суммой списани / зачисления
		if accTransaction.Indicator != "E" {
			err = tx.Model(&account).Update(col, gorm.Expr(colSql, accTransaction.Bonus)).Error
		}

		if err != nil {
			return err
		}

		// Если пролучилось - надо записать транзакцию в журнал
		err = tx.WithContext(ctx).Save(accTransaction).Error
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}
	return terr
}

// Получаем текущий баланс пользователя
func (db *GormDB) GetAccount(ctx context.Context, account *models.Account) error {
	return db.DB.WithContext(ctx).First(account).Error
}

func (db *GormDB) GetHistory(ctx context.Context, transaction *models.Transaction) []models.Transaction {
	var transactions []models.Transaction
	db.DB.WithContext(ctx).Order("time_stamp DESC").Find(&transactions, transaction)
	return transactions
}
