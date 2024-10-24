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

	//TODO пересмотреть интерфейсы ↓↓↓
	PostTransaction(ctx context.Context, transaction *Transaction) error //проводка по счету
}

type AccountStorage interface {
	//создает счет для нового пользователя
	//происходит при регистраии пользователя
	CreateAccount(ctx context.Context, account *Account) error
	//возвращает баланс пользователя
	GetAccount(ctx context.Context, account *Account) error
	//возвращает историю операций
	GetHistory(ctx context.Context, transaction *Transaction) []Transaction

	//TODO пересмотреть интерфейсы ↓↓↓
	PostTransaction(ctx context.Context, transaction *Transaction) error //проводка по счету
}

type UserStorage interface {
	GetUser(ctx context.Context, user *User) error        //возвращает пользователя по id
	CreateUser(ctx context.Context, user *User) error     //создает пользователя
	CheckUserExist(ctx context.Context, user *User) error //проверяем есть ли пользователь с таким имененм
}
