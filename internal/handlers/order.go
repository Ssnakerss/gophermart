package handlers

import (
	"context"
	"io"
	"log/slog"
	"net/http"

	"github.com/Ssnakerss/gophermart/internal/db"
	"github.com/Ssnakerss/gophermart/internal/order"
	"github.com/Ssnakerss/gophermart/internal/user"
)

func PostAPIUserOrders(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Warn("request body read error")
		http.Error(w, "request body read error", http.StatusInternalServerError)
		return
	}

	//базовые проверки пройдены, создает заказ
	op := order.NewOrderProcessor(db.New(db.ConString))
	order := op.NewOrder(context.TODO(), string(body), &user.User{ID: "dummy", IsAuthorizad: true})

	w.WriteHeader(order.Status.ResponseCode())
	w.Write([]byte(order.Status))
}
