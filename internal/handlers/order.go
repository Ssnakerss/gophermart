package handlers

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

func (hm *HandlerMaster) PostAPIUserOrders(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Warn("request body read error")
		http.Error(w, "request body read error", http.StatusInternalServerError)
		return
	}
	//базовые проверки пройдены, создает заказ
	order := hm.orderProcessoe.NewOrder(hm.rootAppContext, string(body), hm.currentUser)
	w.WriteHeader(order.Status.ResponseCode())
	w.Write([]byte(order.Status))
}

func (hm *HandlerMaster) GetAPIUserOrders(w http.ResponseWriter, r *http.Request) {
	//получаем список ордеров по текущему пользователю
	orders := hm.orderProcessoe.AllOrders(hm.rootAppContext, hm.currentUser)
	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent) //204
		return
	}

	body, err := json.Marshal(orders)
	if err != nil {
		slog.Error("GetAPIUserOrders", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError) //500
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
