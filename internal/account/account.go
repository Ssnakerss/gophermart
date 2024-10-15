package account

import (
	"context"

	"github.com/Ssnakerss/gophermart/internal/models"
)

type AccountKeeper struct {
	storage models.AccountStorage
}

func NewAccountKeeper(st models.AccountStorage) *AccountKeeper {
	var ac AccountKeeper
	ac.storage = st
	return &ac
}

func (ac *AccountKeeper) GetAccount(ctx context.Context, id string) (*models.Account, error) {
	acc := models.Account{UserID: id}
	err := ac.storage.GetAccount(ctx, &acc)
	return &acc, err
}

func (ac *AccountKeeper) PostTransaction(ctx context.Context, tr *models.Transaction) error {
	return ac.storage.PostTransaction(context.TODO(), tr)
}

func (ac *AccountKeeper) GetWithdrawHistory(ctx context.Context, userID string) []models.Transaction {
	tr := models.Transaction{
		UserID:    userID,
		Indicator: "C",
	}
	return ac.storage.GetHistory(ctx, &tr)

}
