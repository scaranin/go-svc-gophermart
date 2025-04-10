package routerAPI

import (
	"go-svc-gophermart/internal/handlers"
	"go-svc-gophermart/internal/middlewares"

	"github.com/go-chi/chi"
)

// Формирование ApiRoutes
//
// Входные параметры: h *handlers.URLHandler - структура с конфигурациями
func NewRouter(h *handlers.URLHandler) *chi.Mux {
	mux := chi.NewRouter()
	mux.Use(middlewares.WithLogging)

	mux.Route("/", func(mux chi.Router) {
		mux.Post("/api/user/register", h.PostUserRegister)
		mux.Post("/api/user/login", h.PostUserLogin)
		mux.Post("/api/user/orders", h.PostUserOrders)
		mux.Post("/api/user/balance/withdraw", h.PostUserBalanceWithDraw)

		mux.Get("/api/user/orders", h.GetUserOrders)
		mux.Get("/api/user/balance", h.GetUserBalance)
		mux.Get("/api/user/withdrawals", h.GetUserBalanceWithDrawals)
	})

	return mux
}
