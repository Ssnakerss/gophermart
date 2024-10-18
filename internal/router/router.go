package router

import (
	"github.com/Ssnakerss/gophermart/internal/handlers"
	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	// r.Post("/api/usesr/register", user.Register)
	// r.Post("/api/user/login", user.Login)
	r.With(middleware.AllowContentType("text/plain")).Post("/api/user/orders", handlers.PostAPIUserOrders)
	r.Get("/api/user/orders", handlers.GetAPIUserOrders)
	// r.Get("/api.user/balance", account.Balance)
	// r.Post("/api/user/balance/withdraw", account.Withdraw)
	// r.Get("/api/user/withdrawals", account.History)
	return r
}
