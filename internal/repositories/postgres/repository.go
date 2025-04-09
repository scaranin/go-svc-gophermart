package repositories

import (
	"go-svc-gophermart/internal/models"
	"go-svc-gophermart/internal/repositories"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RepoDBPostgres struct {
	PGXPool *pgxpool.Pool
}

func NewRepoDBPostgres() (RepoDBPostgres, error) {
	return RepoDBPostgres{}, nil
}

func (repoPG RepoDBPostgres) UserRegister(user models.User) error {
	return nil
}

func (repoPG RepoDBPostgres) UserLogin(user models.User) error {
	return nil
}

func (repoPG RepoDBPostgres) UserOrderCreate(order string) error {
	return nil
}

func (repoPG RepoDBPostgres) GetUserOrders() ([]models.Order, error) {
	var orders []models.Order
	return orders, nil
}

func (repoPG RepoDBPostgres) GetBalanceAccrual() (models.OrderAccrual, error) {
	var OrderAccrual models.OrderAccrual
	return OrderAccrual, nil
}

func (repoPG RepoDBPostgres) GetUserBalance() (models.Balance, error) {
	var Balnce models.Balance
	return Balnce, nil
}

func (repoPG RepoDBPostgres) RequestOrderAccrual() (models.BalanceWithDraw, error) {
	var BalanceWithDraw models.BalanceWithDraw
	return BalanceWithDraw, nil
}

func NewRepository(DSN string) (repositories.GopherMart, error) {
	var repoGM repositories.GopherMart
	return repoGM, nil

}
