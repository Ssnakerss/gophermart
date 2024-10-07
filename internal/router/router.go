package router

import (
	chi "github.com/go-chi/chi/v5"
)

func New() *chi.Mux {
	r := chi.NewRouter()
	// r.Post("api/usesr/register", user.Register)
	// r.Post("api/user/login", user.Login)
	// r.Post("api/user/orders", order.List)
	// r.Get("api.user/balance", account.Balance)
	// r.Post("api/user/balance/withdraw", account.Withdraw)
	// r.Get("/api/user/withdrawals", account.History)
	return r
}
