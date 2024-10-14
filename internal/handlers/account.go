package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/Ssnakerss/gophermart/internal/account"
	"github.com/Ssnakerss/gophermart/internal/db"
	"github.com/Ssnakerss/gophermart/internal/models"
	"github.com/Ssnakerss/gophermart/internal/types"
)

func GetAPIUserBalance(w http.ResponseWriter, r *http.Request) {
	ac := account.NewAccountKeeper(db.New(db.ConString, db.Info))
	// currentUser := string(r.Context().Value("UserID"))
	acc, err := ac.GetAccount(context.TODO(), "ivan")
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
	w.Write(body)
}

func PostAPIUserBalanceWithdraw(w http.ResponseWriter, r *http.Request) {
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

	tr.Indicator = "C" //withdrawal
	tr.UserID = "ivan" //current user
	tr.TimeStamp = types.TimeRFC3339(time.Now())
	ac := account.NewAccountKeeper(db.New(db.ConString, db.Info))
	err = ac.PostTransaction(context.TODO(), &tr)

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

func GetAPIUserWithdrawals(w http.ResponseWriter, r *http.Request) {
	ac := account.NewAccountKeeper(db.New(db.ConString, db.Info))
	trs := ac.GetWithdrawHistory(context.TODO(), "ivan") //current user
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

	w.Write(body)
}
