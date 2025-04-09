package models

// Структура JSON. Пользователь
type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
