package handlers

import (
	"go-svc-gophermart/internal/models"
	"log"
	"net/http"
)

// Получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях
func (h *URLHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(h.GetCurrentMethodName()))
}

// Загрузка пользователем номера заказа для расчёта
//
// Хендлер доступен только аутентифицированным пользователям. Номером заказа является последовательность цифр произвольной длины.
// Возможные коды ответа:
//
// - `200` — номер заказа уже был загружен этим пользователем;
// - `202` — новый номер заказа принят в обработку;
// - `400` — неверный формат запроса;
// - `401` — пользователь не аутентифицирован;
// - `409` — номер заказа уже был загружен другим пользователем;
// - `422` — неверный формат номера заказа;
// - `500` — внутренняя ошибка сервера.
func (h *URLHandler) PostUserOrders(w http.ResponseWriter, r *http.Request) {
	err := h.Repo.UserOrderCreate(models.OrderShort{User: "user1", OrderNumber: "order1"})
	if err != nil {
		log.Println(err)
	}
	w.Write([]byte(h.GetCurrentMethodName()))
}
