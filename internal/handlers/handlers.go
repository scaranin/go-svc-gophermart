package handlers

import (
	"go-svc-gophermart/internal/client"
	"go-svc-gophermart/internal/middlewares"
	"go-svc-gophermart/internal/repositories"
	"log"
	"net/http"
	"runtime"
	"unicode"
)

// Основная структура со списком обработчиков
type URLHandler struct {
	Repo       repositories.GopherMart
	TokenSvc   middlewares.TokenService
	AccrualSvc client.AccrualService
}

// Функция для получения имени текущего метода
func (h URLHandler) GetCurrentMethodName() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return "unknown"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}

	fullName := fn.Name()

	return fullName
}

// Обработка cookie
func (h *URLHandler) cookieProcessing(w http.ResponseWriter, r *http.Request) (string, error) {
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

// Проверка номера договора алгоритмом Луна
func validateLuhn(orderNumber string) bool {
	for _, r := range orderNumber {
		if !unicode.IsDigit(r) {
			return false
		}
	}

	sum := 0
	for i, rune := range orderNumber {
		digit := int(rune)
		if (len(orderNumber)-i)%2 == 0 {
			digit *= 2

			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}

	return sum%10 == 0
}
