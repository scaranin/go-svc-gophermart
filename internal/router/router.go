package router

import (
	"go-svc-gophermart/internal/handlers"
	"go-svc-gophermart/internal/middlewares"

	"github.com/go-chi/chi"
)

// Формирование Routes
//
// Входные параметры: h *handlers.URLHandler - структура с конфигурациями
func NewRouter(h *handlers.URLHandler, authConfig middlewares.AuthConfig) *chi.Mux {
	mux := chi.NewRouter()
	authService := middlewares.NewAuthService(authConfig.SecretKey)

	authMiddleware := middlewares.WithAuth(&authService)

	mux.Use(middlewares.WithLogging, authMiddleware)

	mux.Route("/", func(mux chi.Router) {
		mux.Post("/api/user/register", h.PostUserRegister)
		mux.Post("/api/user/login", h.PostUserLogin)
		mux.Post("/api/user/orders", h.PostUserOrders)
		mux.Post("/api/user/balance/withdraw", h.RequestWithdraw)

		mux.Get("/api/user/orders", h.GetUserOrders)
		mux.Get("/api/user/balance", h.GetUserBalance)
		mux.Get("/api/user/withdrawals", h.GetWithdrawals)
	})

	return mux
}
