package db

import (
	"context"

	"github.com/Ssnakerss/gophermart/internal/apperrs"
	"github.com/Ssnakerss/gophermart/internal/models"
	"gorm.io/gorm"
)

// Account related operations
// Сначала пробуем обновить баланс -  важно при списании
// делаем все в транзакции для сохранения целостности данных
func (db *GormDB) PostTransaction(ctx context.Context, accTransaction *models.Transaction) error {

	if accTransaction.Indicator != "D" && accTransaction.Indicator != "C" {
		return apperrs.ErrInvalidIndicator
	}
	account := models.Account{
		UserID: accTransaction.UserID,
	}

	col := ""
	colSQL := ""

	var terr error
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var err error
		switch accTransaction.Indicator {
		//пополнение счета
		case "D":
			err = tx.Model(&account).Update("balance", gorm.Expr("balance + ?", accTransaction.Accrual)).Error
			if err != nil {
				return err
			}
			col = "debit"
			colSQL = "debit + ?"
		//списание со счета
		case "C":
			result := tx.Model(&account).
				//при списании надо проверить достаочно ли средств на счету
				Where("balance >= ?", accTransaction.Accrual).
				Update("balance", gorm.Expr("balance - ?", accTransaction.Accrual))
				//если списание не прощло - ставим ошибку
			if result.RowsAffected == 0 {
				terr = apperrs.ErrInsufficientFunds
			} else {
				err = result.Error
			}

			col = "credit"
			colSQL = "credit + ?"
		}

		if terr == apperrs.ErrInsufficientFunds {
			accTransaction.Indicator = "E"
		}

		if err != nil {
			return err
		}
		//если транзакция обновила баланс успешно - обноаляем и пола счета с суммой списани / зачисления
		if accTransaction.Indicator != "E" {
			err = tx.Model(&account).Update(col, gorm.Expr(colSQL, accTransaction.Accrual)).Error
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

// создание счетя для новых пользователей
func (db *GormDB) CreateAccount(ctx context.Context, account *models.Account) error {
	return db.DB.WithContext(ctx).Create(account).Error
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
