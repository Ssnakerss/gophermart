package db

import (
	"context"

	"github.com/Ssnakerss/gophermart/internal/models"
	"gorm.io/gorm"
)

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
