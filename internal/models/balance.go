package models

import "time"

// Структура JSON. Получение текущего баланса пользователя
type Balance struct {
	Current   float32 `json:"current"`
	WithDrawn float32 `json:"withdrawn"`
}

// Структура JSON. Запрос на списание средств
type RequestWithDraw struct {
	Order string  `json:"order"`
	Sum   float32 `json:"sum"`
}

// Структура JSON. Получение информации о выводе средств
type WithDrawalsList struct {
	Order       string    `json:"order"`
	Sum         float32   `json:"sum"`
	DtProcessed time.Time `json:"processed_at"` //Формат даты RFC3339
}
