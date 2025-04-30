package handlers

import (
	"go-svc-gophermart/internal/client"
	"go-svc-gophermart/internal/middlewares"
	"go-svc-gophermart/internal/repositories"
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
