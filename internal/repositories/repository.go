package repositories

import (
	"go-svc-gophermart/internal/models"
)

type User interface {
	UserRegister(user models.User) error
	UserLogin(user models.User) error
}

type Order interface {
	UserOrderCreate(order models.OrderShort) error
	GetAccruals(User string) ([]models.OrderAccrual, error)
	UpdateOrderList(OrderAccrualArr []models.OrderAccrual) error
	GetUserOrders(User string) ([]models.Order, error)
}

type Balance interface {
	GetUserBalance(User string) (models.Balance, error)
	RequestOrderAccrual(models.BalanceWithDraw) error
}

type GopherMart interface {
	User
	Order
	Balance
}
