package models

import (
	"database/sql"
	"time"
)

// Структура JSON. Получение текущего баланса пользователя
type Balance struct {
	Current   sql.NullFloat64 `json:"current"`
	WithDrawn sql.NullFloat64 `json:"withdrawn"`
}

// Структура JSON. Запрос на списание средств
type RequestWithDraw struct {
	Order string          `json:"order"`
	Sum   sql.NullFloat64 `json:"sum"`
}

// Структура JSON. Получение информации о выводе средств
type WithDrawalsList struct {
	Order       string          `json:"order"`
	Sum         sql.NullFloat64 `json:"sum"`
	DtProcessed time.Time       `json:"processed_at"` //Формат даты RFC3339
}
