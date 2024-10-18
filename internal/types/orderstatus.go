package types

// --------------------------------------------------------------------------------
// статус заказа с сопоставлением кода ответа
type OrderStatus string

// статусы по итогам обработки нового заказа
const (
	CHECKING OrderStatus = "CHECKING" //проверка в процессе

	REPEATED    OrderStatus = "REPEATED"    //заказ уже был загружен этим пользователем
	DUPLICATED  OrderStatus = "DUPLICATED"  //заказ уже был загружен другим пользователем
	WRONGFORMAT OrderStatus = "WRONGFORMAT" //неправильный формат номера заказа

	ERROR OrderStatus = "ERROR"
)

// Статусы принятых заказаов
const (
	//промежуточные статусы
	NEW        OrderStatus = "NEW"        //заказ загружен но еще не приянт в обработку
	REGISTERED OrderStatus = "REGISTERED" //заказ зарегистрирован но бонус не расчитан
	PROCESSING OrderStatus = "PROCESSING" //заказ в оборботке бонусы расчитыватся

	//финальлные статусы
	INVALID OrderStatus = "INVALID" //система расчета бонусов отказала в расчете

	//такие заказы пишем и в историю по счету
	//Debit-Credit indicator = D
	PROCESSED OrderStatus = "PROCESSED" //заказ обработан, бонусы получены - пишем в дебет счета
)

// статуся для операций с  бонусами
// заказы с этими статусами пишем только в историю по счету
const (
	//это списание бонусов по заказу - пишем в кредит счета
	//Debit-Credit indicator = C
	WITHDRAW OrderStatus = "WITHDRAW"
	//бонусов не хватает списания - для истории
	//Debit-Credit indicator = E
	NOTENOUGHBONUS OrderStatus = "NOTENOUGHBONUS"
)

func (os OrderStatus) ResponseCode() int {
	switch os {
	case REPEATED:
		return 200
	case NEW:
		return 202
	case DUPLICATED:
		return 409
	case WRONGFORMAT:
		return 422
	case ERROR:
		return 500
	default:
		return 0
	}
}
