package models

import "time"

// Базовая структура заказа
type OrderShort struct {
	User        string
	OrderNumber string
}

// Структура JSON. Получение списка загруженных номеров заказов
type Order struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual"`
	DtUploaded time.Time `json:"uploaded_at"` //Формат даты RFC3339
}

// Структура JSON. Получение информации о расчёте начислений баллов лояльности
type OrderAccrual struct {
	Order   string `json:"order"`
	Status  string `json:"status"`
	Accrual int    `json:"accrual"`
}
