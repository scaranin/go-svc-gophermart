package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-svc-gophermart/internal/models"
	"log"
	"net/http"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func (h *URLHandler) CookieProcessing(w http.ResponseWriter, r *http.Request) (string, error) {
	var (
		user string
		err  error
	)
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return user, err
	}
	user, err = h.TokenSvc.GetUserFromCookie(cookie)
	if err != nil {
		log.Println(err)
	}

	cookieW, err := h.TokenSvc.GenerateCookie(user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
	}
	http.SetCookie(w, cookieW)

	return user, err
}

// Получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях
//
// Возможные коды ответа:
//
// - `200` — успешная обработка запроса.
// - `204` — нет данных для ответа.
// - `401` — пользователь не авторизован.
// - `500` — внутренняя ошибка сервера.
func (h *URLHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := h.CookieProcessing(w, r)

	if err != nil {
		log.Println(err)
		return
	}

	// Получаем список заказов пользователя не в конечном статусе
	OrderAccrualList, err := h.Repo.GetAccruals(user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Println("OrderAccrual ", OrderAccrualList)
	var OrderList []models.OrderAccrual
	for _, OrderItem := range OrderAccrualList {
		Order, err := h.AccrualSvc.GetOrder(OrderItem.Order)
		if err != nil {
			log.Println(err)
			continue
		}
		OrderList = append(OrderList, Order)

	}
	// Обновляем
	if len(OrderList) != 0 {
		if err = h.Repo.UpdateOrderList(OrderList); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	// Получаем заказы с актуальным статусом
	Orders, err := h.Repo.GetUserOrders(user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	OrdersJSON, err := json.Marshal(Orders)
	if err != nil {
		log.Println(err)

		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(OrdersJSON)
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
	var (
		buf    bytes.Buffer
		header int
		pgErr  *pgconn.PgError
		ok     bool
	)

	user, err := h.CookieProcessing(w, r)

	if err != nil {
		log.Println(err)
		return
	}

	defer r.Body.Close()
	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	orderNumber := buf.String()

	if !validateLuhn(orderNumber) {
		pgErr, ok = h.Repo.UserOrderCreate(models.OrderShort{User: user, OrderNumber: orderNumber}).(*pgconn.PgError)
	} else {
		header = http.StatusUnprocessableEntity
		log.Println("Bad order number")
	}

	if ok {
		switch pgErr.Code {
		case pgerrcode.SuccessfulCompletion:
			header = http.StatusAccepted
			log.Println(pgErr)
		case "-1":
			header = http.StatusOK
			log.Println(pgErr)
		case "-2":
			header = http.StatusConflict
			log.Println(pgErr)
		default:
			header = http.StatusInternalServerError
			log.Print(pgErr)
		}
	}

	if header == 0 {
		header = http.StatusAccepted
	}
	w.WriteHeader(header)
}
