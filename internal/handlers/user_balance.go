package handlers

import (
	"net/http"
)

// Получение текущего баланса счёта баллов лояльности пользователя
func (h *URLHandler) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(h.GetCurrentMethodName()))
}

// Запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа
func (h *URLHandler) PostUserBalanceWithDraw(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(h.GetCurrentMethodName()))
}

// Получение информации о выводе средств с накопительного счёта пользователем
func (h *URLHandler) GetUserBalanceWithDrawals(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(h.GetCurrentMethodName()))
}
