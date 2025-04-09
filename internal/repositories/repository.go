package repositories

import (
	"go-svc-gophermart/internal/models"
)

type User interface {
	UserRegister(user models.User) error
	UserLogin(user models.User) error
}

type Order interface {
	UserOrderCreate(order string) error
	GetUserOrders() ([]models.Order, error)
	GetBalanceAccrual() (models.OrderAccrual, error)
}

type Balance interface {
	GetUserBalance() (models.Balance, error)
	RequestOrderAccrual(models.BalanceWithDraw, error) error
}

type GopherMart interface {
	User
	Order
	Balance
}

type Repository interface {
	NewRepository(DSN string) (GopherMart, error)
}
