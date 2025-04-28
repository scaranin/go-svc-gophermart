package handlers

import (
	"bytes"
	"fmt"
	"go-svc-gophermart/internal/models"
	"log"
	"net/http"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

// Получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях
//
// Возможные коды ответа:
//
// - `200` — успешная обработка запроса.
// - `204` — нет данных для ответа.
// - `401` — пользователь не авторизован.
// - `500` — внутренняя ошибка сервера.
func (h *URLHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user, err := h.TokenSvc.GetUserFromCookie(cookie)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("user ", user)

	// Получаем список заказов пользователя не в конечном статусе и обновляем данные в БД
	h.Repo.GetUserOrders(user, []string{"INVALID", "PROCESSED"})
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
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user, err := h.TokenSvc.GetUserFromCookie(cookie)
	if err != nil {
		log.Println(err)
	}
	fmt.Print("user", user)

	var buf bytes.Buffer
	defer r.Body.Close()
	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	orderNumber := buf.String()
	fmt.Print("orderNumber", orderNumber)

	var header int
	pgErr, ok := h.Repo.UserOrderCreate(models.OrderShort{User: user, OrderNumber: orderNumber}).(*pgconn.PgError)

	if ok {
		switch pgErr.Code {
		case pgerrcode.SuccessfulCompletion:
			header = http.StatusAccepted
			log.Print(pgErr)
		case "-1":
			header = http.StatusOK
			log.Print(pgErr)
		case "-2":
			header = http.StatusConflict
			log.Print(pgErr)
		default:
			header = http.StatusInternalServerError
			log.Print(pgErr)
		}
	}

	cookieW, err := h.TokenSvc.GenerateCookie(user)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Print(err.Error())
	}
	http.SetCookie(w, cookieW)
	if header == 0 {
		header = http.StatusAccepted
	}
	w.WriteHeader(header)
}
