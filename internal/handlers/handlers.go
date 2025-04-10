package handlers

import (
	"go-svc-gophermart/internal/repositories"
	"runtime"
)

// Основная структура со списком обработчиков
type URLHandler struct {
	Repo repositories.GopherMart
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
