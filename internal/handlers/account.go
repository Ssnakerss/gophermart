package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/Ssnakerss/gophermart/internal/models"
	"github.com/Ssnakerss/gophermart/internal/types"
)

func (hm *HandlerMaster) GetAPIUserBalance(w http.ResponseWriter, r *http.Request) {

	acc, err := hm.accountKeeper.GetAccount(hm.rootAppContext, hm.currentUser.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error(err.Error())
		return
	}
	body, err := json.Marshal(acc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error(err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) //200
	w.Write(body)
}

func (hm *HandlerMaster) PostAPIUserBalanceWithdraw(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Warn("request body read error")
		http.Error(w, "request body read error", http.StatusInternalServerError)
		return
	}
	tr := models.Transaction{}
	err = json.Unmarshal(body, &tr)
	if err != nil {
		slog.Warn(err.Error())
		http.Error(w, "transaction parse error", http.StatusUnprocessableEntity)
		return
	}
	//проверяем правильность номера заказа
	if err = tr.OrderNumber.Check(); err != nil {
		http.Error(w, "order number is invalid", http.StatusUnprocessableEntity)
		slog.Warn(err.Error())
		return
	}

	tr.Indicator = "C"            //withdrawal
	tr.UserID = hm.currentUser.ID //current user
	tr.TimeStamp = types.TimeRFC3339(time.Now())

	err = hm.accountKeeper.PostTransaction(hm.rootAppContext, &tr)
	if err != nil {
		if errors.Is(err, models.ErrInsufficientFunds) {
			http.Error(w, "insufficient funds", http.StatusPaymentRequired)
			slog.Warn(err.Error())
			return
		}
		http.Error(w, "transaction error", http.StatusInternalServerError)
		slog.Error(err.Error())
		return
	}
}

func (hm *HandlerMaster) GetAPIUserWithdrawals(w http.ResponseWriter, r *http.Request) {

	trs := hm.accountKeeper.GetWithdrawHistory(hm.rootAppContext, hm.currentUser.ID) //current user
	if len(trs) == 0 {
		http.Error(w, "no withdrawals", http.StatusNoContent)
		slog.Warn("no withdrawals")
		return
	}

	body, err := json.Marshal(trs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error(err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
