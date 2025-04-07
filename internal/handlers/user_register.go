package handlers

import (
	"net/http"
)

// Регистрация пользователя
func (h *URLHandler) PostUserRegister(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(h.GetCurrentMethodName()))
}
