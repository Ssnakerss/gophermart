package router

import (
	"github.com/Ssnakerss/gophermart/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

func New(hm *handlers.HandlerMaster) *chi.Mux {
	r := chi.NewRouter()
	// chi middleware
	r.Use(middleware.Logger)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AllowContentType("application/json"))
		//авторизация и регистрация пользователя
		r.Post("/api/user/register", hm.PostAPIUserRegister)
		r.Post("/api/user/login", hm.PostAPIUserLogin)
	})

	//ручки доступные только после авторизации
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(hm.UserManager.GetTokenAuth()))
		r.Use(jwtauth.Authenticator(hm.UserManager.GetTokenAuth()))
		r.Use(hm.AuthorizeUser)

		//работа с заказамт
		r.With(middleware.AllowContentType("text/plain")).Post("/api/user/orders", hm.PostAPIUserOrders)
		r.Get("/api/user/orders", hm.GetAPIUserOrders)
		//работа с счетами пользователя
		r.Get("/api/user/balance", hm.GetAPIUserBalance)
		r.With(middleware.AllowContentType("application/json")).Post("/api/user/balance/withdraw", hm.PostAPIUserBalanceWithdraw)
		r.Get("/api/user/withdrawals", hm.GetAPIUserWithdrawals)
	})

	return r
}
