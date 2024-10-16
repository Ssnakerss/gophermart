package models

import (
	"context"
)

type OrderStorage interface {
	//вставлят или обновляет заказ в хранилище по номеру
	SaveOrder(ctx context.Context, order *Order) error
	//достает заказ из хранилища по номеру
	GetOrder(ctx context.Context, order *Order) error
	//достает все заказы из хранилища по параметрам
	//в нашем случае это имя пользователя и/или статус
	GetAllOrders(ctx context.Context, order *Order) []Order

	GetOrdersByStatus(ctx context.Context, order []Order) []Order //возвращает заказы по статусу
}

type AccountStorage interface {
	CreateAccount(ctx context.Context, account *Account) error
	GetAccount(ctx context.Context, account *Account) error                 //возвращает баланс пользователя
	PostTransaction(ctx context.Context, transaction *Transaction) error    //проводка по счету
	GetHistory(ctx context.Context, transaction *Transaction) []Transaction //возвращает историю операций
}

type UserStorage interface {
	GetUser(ctx context.Context, user *User) error        //возвращает пользователя по id
	CreateUser(ctx context.Context, user *User) error     //создает пользователя
	CheckUserExist(ctx context.Context, user *User) error //проверяем есть ли пользователь с таким имененм
}
