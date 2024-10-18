package handlers

import (
	"context"

	"github.com/Ssnakerss/gophermart/internal/db"
	"github.com/Ssnakerss/gophermart/internal/models"
	"github.com/Ssnakerss/gophermart/internal/proc/account"
	"github.com/Ssnakerss/gophermart/internal/proc/order"
	"github.com/Ssnakerss/gophermart/internal/proc/user"
)

type HandlerMaster struct {
	orderProcessoe *order.OrderProcessor
	accountKeeper  *account.AccountKeeper
	UserManager    *user.UserManager

	currentUser *models.User

	storage *db.GormDB

	rootAppContext context.Context
}

func NewMaster(ctx context.Context, st *db.GormDB) *HandlerMaster {
	var hm HandlerMaster
	hm.storage = st
	hm.orderProcessoe = order.NewOrderProcessor(st)
	hm.accountKeeper = account.NewAccountKeeper(st)
	hm.UserManager = user.NewUserManager(st)

	hm.rootAppContext = ctx //TODO	- ???

	return &hm
}
