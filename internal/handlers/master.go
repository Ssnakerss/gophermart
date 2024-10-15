package handlers

import (
	"context"

	"github.com/Ssnakerss/gophermart/internal/account"
	"github.com/Ssnakerss/gophermart/internal/db"
	"github.com/Ssnakerss/gophermart/internal/models"
	"github.com/Ssnakerss/gophermart/internal/order"
	"github.com/Ssnakerss/gophermart/internal/user"
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

//TODO - implement

// func (hm *HandlerMaster) GetTokenAuth() *jwtauth.JWTAuth {
// 	return hm.userManager.GetTokenAuth()
// }

// func (hm HandlerMaster) GetUserChecker() http.Handler(
// 	return userManager.AuthorizeUser
// )
