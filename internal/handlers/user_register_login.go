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

type RegisterOrLogin func(models.User) error

// Регистрация или аутентификация пользователя.
func (h *URLHandler) UserRegisterOrLogin(w http.ResponseWriter, r *http.Request, t string) {
	var (
		err    error
		user   models.User
		resp   []byte
		buf    bytes.Buffer
		header int
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

	var registerOrLogin RegisterOrLogin

	if t == "Register" {
		registerOrLogin = h.Repo.UserRegister
	}

	if t == "Login" {
		registerOrLogin = h.Repo.UserLogin
	}

	pgErr, ok := registerOrLogin(user).(*pgconn.PgError)
	if ok {
		switch pgErr.Code {
		case pgerrcode.SuccessfulCompletion:
			header = http.StatusOK
			log.Print(pgErr)
		case pgerrcode.UniqueViolation:
			header = http.StatusConflict
			log.Print(pgErr)
		case pgerrcode.InvalidAuthorizationSpecification:
			header = http.StatusUnauthorized
			log.Print(pgErr)
		default:
			header = http.StatusInternalServerError
			log.Print(pgErr)
		}
	}

	cookieW, err := h.Auth.FillUserCookie(user.Login)
	if err != nil {
		log.Print(err.Error())
	}
	http.SetCookie(w, cookieW)
	if header == 0 {
		header = http.StatusOK
	}
	w.WriteHeader(header)
	w.Write(resp)
}

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
	h.UserRegisterOrLogin(w, r, "Register")
}

// Аутентификация пользователя
//
// Аутентификация производится по паре логин/пароль.\
//
// Возможные коды ответа:
//
// - `200` — пользователь успешно зарегистрирован и аутентифицирован;
// - `400` — неверный формат запроса;
// - `401` — неверная пара логин/пароль;
// - `500` — внутренняя ошибка сервера.
func (h *URLHandler) PostUserLogin(w http.ResponseWriter, r *http.Request) {
	h.UserRegisterOrLogin(w, r, "Login")
}
