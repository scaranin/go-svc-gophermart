package handlers

import (
	"bytes"
	"encoding/json"
	"go-svc-gophermart/internal/models"
	"log"
	"net/http"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

// Регистрация пользователя
//
// Регистрация производится по паре логин/пароль.
// После успешной регистрации выполняется автоматическая аутентификация пользователя.
//
// Возможные коды ответа:
//
// - `200` — пользователь успешно зарегистрирован и аутентифицирован;
// - `400` — неверный формат запроса;
// - `409` — логин уже занят;
// - `500` — внутренняя ошибка сервера.
func (h *URLHandler) PostUserRegister(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		user models.User
		resp []byte
		buf  bytes.Buffer
	)

	defer r.Body.Close()
	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pgErr, ok := h.Repo.UserRegister(user).(*pgconn.PgError)
	if ok {
		if pgErr.Code == pgerrcode.UniqueViolation {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(pgErr)
		}
	}

	cookieW, err := h.Auth.FillUserCookie(user.Login)
	if err != nil {
		log.Print(err.Error())
	}
	http.SetCookie(w, cookieW)

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
