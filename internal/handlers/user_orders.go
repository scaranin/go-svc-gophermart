package handlers

import (
	"net/http"
)

// Получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях
func (h *URLHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(h.GetCurrentMethodName()))
}

// Загрузка пользователем номера заказа для расчёта
func (h *URLHandler) PostUserOrders(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(h.GetCurrentMethodName()))
}
