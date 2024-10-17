package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/Ssnakerss/gophermart/internal/models"
)

func (hm *HandlerMaster) PostAPIUserLogin(w http.ResponseWriter, r *http.Request) {
	//достаем логин пароль из боди
	cr, err := getCred(r.Body)
	if err != nil {
		slog.Warn("login error", "credentials", cr, "error", err.Error())
		if errors.Is(err, models.ErrStatusBadRequest) {
			http.Error(w, "body parse error", http.StatusBadRequest)
			return
		}
		http.Error(w, "request body read error", http.StatusInternalServerError)
		return
	}

	hm.currentUser, err = hm.UserManager.Login(hm.rootAppContext, cr)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			slog.Warn("login erorr", "user not found", err.Error())
			http.Error(w, "user not found", http.StatusUnauthorized)
			return
		}
		slog.Warn("login error", " internal server", err)
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
	w.WriteHeader(http.StatusOK)
	slog.Info("user login success", "ID", hm.currentUser.ID)
}

func (hm *HandlerMaster) PostAPIUserRegister(w http.ResponseWriter, r *http.Request) {
	//достаем логин пароль из боди
	cr, err := getCred(r.Body)
	if err != nil {
		slog.Warn("register error", "credentials", cr, "error", err.Error())
		if errors.Is(err, models.ErrStatusBadRequest) {
			http.Error(w, "body parse error", http.StatusBadRequest)
			return
		}
		http.Error(w, "request body read error", http.StatusInternalServerError)
		return
	}

	hm.currentUser, err = hm.UserManager.Register(hm.rootAppContext, cr)
	if err != nil {
		if errors.Is(err, models.ErrUserAlreadyExists) {
			slog.Warn("register error", "user already exist", err.Error())
			http.Error(w, "user already exist", http.StatusConflict) //409
			return
		}
		slog.Warn("register", " error", err)
		http.Error(w, "user register error", http.StatusInternalServerError)
		return
	}
	//TODO - сделать создание пользователя и счета в одной транзакции
	//создадим счет для нового пользователя
	_, err = hm.accountKeeper.CreateAccount(hm.rootAppContext, hm.currentUser.ID)
	if err != nil {
		slog.Warn("account create", " error", err)
		http.Error(w, "account create error", http.StatusInternalServerError)
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
	w.WriteHeader(http.StatusOK)
	slog.Info("user register success", "ID", hm.currentUser.ID)
}

func getCred(rBody io.ReadCloser) (*models.UserCred, error) {
	body, err := io.ReadAll(rBody)
	if err != nil {
		return nil, err
	}

	cr := models.UserCred{}
	err = json.Unmarshal(body, &cr)
	if err != nil {
		return nil, fmt.Errorf("%w:%w", models.ErrStatusBadRequest, err)
	}
	return &cr, nil
}
