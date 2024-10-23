package account

import (
	"context"
	"testing"

	"github.com/Ssnakerss/gophermart/internal/mock"
	"github.com/Ssnakerss/gophermart/internal/models"
)

func TestAccountKeeper_PostTransaction(t *testing.T) {

	mockAccountStorage := mock.NewMockAccountStorage()
	ac := NewAccountKeeper(mockAccountStorage)
	ctx := context.Background()
	ac.CreateAccount(ctx, "test")

	type args struct {
		ctx context.Context
		tr  *models.Transaction
	}
	tests := []struct {
		name    string
		ac      *AccountKeeper
		args    args
		wantErr bool
	}{
		{
			name: "Debet",
			ac:   ac,
			args: args{
				ctx: ctx,
				tr: &models.Transaction{
					UserID:      "test",
					OrderNumber: 1000256,
					Accrual:     1000,
					Indicator:   "D",
				},
			},
			wantErr: false,
		},
		{
			name: "Credi OK",
			ac:   ac,
			args: args{
				ctx: ctx,
				tr: &models.Transaction{
					UserID:      "test",
					OrderNumber: 1000256,
					Accrual:     1000,
					Indicator:   "C",
				},
			},
			wantErr: false,
		},
		{
			name: "Credit NG",
			ac:   ac,
			args: args{
				ctx: ctx,
				tr: &models.Transaction{
					UserID:      "test",
					OrderNumber: 1000256,
					Accrual:     1000,
					Indicator:   "C",
				},
			},
			wantErr: true,
		},
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.ac.PostTransaction(tt.args.ctx, tt.args.tr); (err != nil) != tt.wantErr {
				t.Errorf("AccountKeeper.PostTransaction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
