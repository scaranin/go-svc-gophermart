package models

import "time"

// Структура JSON. Получение текущего баланса пользователя
type Balance struct {
	Current   int `json:"current"`
	WithDrawn int `json:"withdrawn"`
}

// Структура JSON. Запрос на списание средств
type BalanceWithDraw struct {
	Order string  `json:"order"`
	Sum   float32 `json:"sum"`
}

// Структура JSON. Получение информации о выводе средств
type BalanceWithDrawals struct {
	Order       string    `json:"order"`
	Sum         float32   `json:"sum"`
	DtProcessed time.Time `json:"processed_at"` //Формат даты RFC3339
}
