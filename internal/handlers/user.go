package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/Ssnakerss/gophermart/internal/models"
)

func (hm *HandlerMaster) PostAPIUserLogin(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Warn("request body read error")
		http.Error(w, "request body read error", http.StatusInternalServerError)
		return
	}

	cr := models.UserCred{}
	err = json.Unmarshal(body, &cr)
	if err != nil {
		slog.Warn(err.Error())
		http.Error(w, "body parse error", http.StatusBadRequest)
		return
	}

	hm.currentUser, err = hm.UserManager.Login(hm.rootAppContext, &cr)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			slog.Warn("user not found", "err", err.Error())
			http.Error(w, "user not found", http.StatusUnauthorized)
			return
		}
		slog.Warn("login", " error", err)
		http.Error(w, "error during login process", http.StatusInternalServerError)
		return
	}

	jwt, err := hm.UserManager.CreateJWT(hm.currentUser)
	if err != nil {
		slog.Warn("jwt create", " error", err)
		http.Error(w, "request body read error", http.StatusInternalServerError)
		return
	}
	// c := http.Cookie{
	// 	Name:  "jwt",
	// 	Value: jwt,
	// }

	// http.SetCookie(w, &c)
	w.Header().Add("Authorization", "Bearer "+jwt)

}

func (hm *HandlerMaster) PostAPIUserRegister(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Warn("request body read error")
		http.Error(w, "request body read error", http.StatusInternalServerError)
		return
	}

	cr := models.UserCred{}
	err = json.Unmarshal(body, &cr)
	if err != nil {
		slog.Warn(err.Error())
		http.Error(w, "body parse error", http.StatusBadRequest)
		return
	}

	hm.currentUser, err = hm.UserManager.Register(hm.rootAppContext, &cr)
	if err != nil {
		if errors.Is(err, models.ErrUserExists) {
			slog.Warn("user already exist", "err", err.Error())
			http.Error(w, "user not found", http.StatusConflict) //409
			return
		}
		slog.Warn("login", " error", err)
		http.Error(w, "error during login process", http.StatusInternalServerError)
		return
	}

	jwt, err := hm.UserManager.CreateJWT(hm.currentUser)
	if err != nil {
		slog.Warn("jwt create", " error", err)
		http.Error(w, "request body read error", http.StatusInternalServerError)
		return
	}
	// c := http.Cookie{
	// 	Name:  "jwt",
	// 	Value: jwt,
	// }

	// http.SetCookie(w, &c)
	w.Header().Add("Authorization", "Bearer "+jwt)
}
