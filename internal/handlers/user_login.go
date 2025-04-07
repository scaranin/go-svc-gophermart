package handlers

import (
	"net/http"
)

// Аутентификация пользователя
func (h *URLHandler) PostUserLogin(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(h.GetCurrentMethodName()))
}
