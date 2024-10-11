package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/Ssnakerss/gophermart/internal/db"
	"github.com/Ssnakerss/gophermart/internal/order"
	"github.com/Ssnakerss/gophermart/internal/user"
)

func PostAPIUserOrders(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Warn("request body read error")
		http.Error(w, "request body read error", http.StatusInternalServerError)
		return
	}
	//базовые проверки пройдены, создает заказ
	op := order.NewOrderProcessor(db.New(db.ConString, db.Info))
	order := op.NewOrder(context.TODO(), string(body), &user.User{ID: "dummy", IsAuthorized: true})

	w.WriteHeader(order.Status.ResponseCode())
	w.Write([]byte(order.Status))
}

func GetAPIUserOrders(w http.ResponseWriter, r *http.Request) {
	//получаем список ордеров по текущему пользователю
	op := order.NewOrderProcessor(db.New(db.ConString, db.Info))
	orders := op.AllOrders(context.TODO(), &user.User{ID: "dummy", IsAuthorized: true})
	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent) //204
		return
	}

	body, err := json.Marshal(orders)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) //500
		return
	}
	w.Write(body)

}
