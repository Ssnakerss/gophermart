package mock

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/Ssnakerss/gophermart/internal/apperrs"
	"github.com/Ssnakerss/gophermart/internal/models"
)

type mockAccountStorage struct {
	accounts     map[string]*models.Account
	transactions map[uint64]*models.Transaction
}

func NewMockAccountStorage() *mockAccountStorage {
	return &mockAccountStorage{
		accounts:     make(map[string]*models.Account),
		transactions: make(map[uint64]*models.Transaction),
	}
}

func (m *mockAccountStorage) CreateAccount(ctx context.Context, account *models.Account) error {
	_, ok := m.accounts[account.UserID]
	if ok {
		return fmt.Errorf("account with id  %s already exists", account.UserID)
	}
	m.accounts[account.UserID] = account
	return nil
}

func (m *mockAccountStorage) GetAccount(ctx context.Context, account *models.Account) error {
	acc, ok := m.accounts[account.UserID]

	if !ok {
		return apperrs.ErrRecordNotFound
	}
	account.Balance = acc.Balance
	account.Debit = acc.Debit
	account.Credit = acc.Credit
	account.UpdatedAt = acc.UpdatedAt
	return nil
}

func (m *mockAccountStorage) GetHistory(ctx context.Context,
	accTransaction *models.Transaction) []models.Transaction {
	res := make([]models.Transaction, 0)
	match := false
	for _, transaction := range m.transactions {
		if accTransaction.UserID == transaction.UserID {
			match = true
			if accTransaction.OrderNumber != 0 {
				match = accTransaction.OrderNumber == transaction.OrderNumber
			}
			if accTransaction.Indicator != "" {
				match = accTransaction.Indicator == transaction.Indicator
			}
			if match {
				res = append(res, *transaction)
			}
		}
	}
	return res
}

func (m *mockAccountStorage) PostTransaction(ctx context.Context, accTransaction *models.Transaction) error {
	if accTransaction.Indicator != "D" && accTransaction.Indicator != "C" {
		return apperrs.ErrInvalidIndicator
	}

	account := models.Account{
		UserID: accTransaction.UserID,
	}

	if err := m.GetAccount(ctx, &account); err != nil {
		return fmt.Errorf("account with id %s not found", accTransaction.UserID)
	}

	if accTransaction.Indicator != "D" && accTransaction.Indicator != "C" {
		return apperrs.ErrInvalidIndicator
	}

	sec1 := rand.New(rand.NewSource(10))
	accTransaction.ID = sec1.Uint64()

	if accTransaction.Indicator == "D" {
		m.accounts[accTransaction.UserID].Debit += accTransaction.Accrual
		m.accounts[accTransaction.UserID].Balance += accTransaction.Accrual
		m.transactions[accTransaction.ID] = accTransaction
	}

	if accTransaction.Indicator == "C" {
		if account.Balance < accTransaction.Accrual {
			return apperrs.ErrInsufficientFunds
		}
		m.accounts[accTransaction.UserID].Credit += accTransaction.Accrual
		m.accounts[accTransaction.UserID].Balance -= accTransaction.Accrual
		m.transactions[accTransaction.ID] = accTransaction
	}

	return nil
}
