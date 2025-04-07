package routerAPI

import (
	"go-svc-gophermart/internal/handlers"
	"go-svc-gophermart/internal/middlewares"

	"github.com/go-chi/chi"
)

/*
	==>

Формирование ApiRoutes
входные параметры: h *handlers.URLHandler - структура с конфигурациями
==<
*/
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

		/*
					* `POST /api/user/register` — регистрация пользователя;
			* `POST /api/user/login` — аутентификация пользователя;
			* `POST /api/user/orders` — загрузка пользователем номера заказа для расчёта;
			* `GET /api/user/orders` — получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
			* `GET /api/user/balance` — получение текущего баланса счёта баллов лояльности пользователя;
			* `POST /api/user/balance/withdraw` — запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
			* `GET /api/user/withdrawals`
		*/

	})

	return mux
}
