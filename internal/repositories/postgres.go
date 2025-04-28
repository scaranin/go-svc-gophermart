package repositories

import (
	"context"
	"fmt"
	"go-svc-gophermart/internal/client"
	"go-svc-gophermart/internal/config"
	"go-svc-gophermart/internal/models"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepoDBPostgres struct {
	PGXPool *pgxpool.Pool
}

func NewRepoDBPostgres() (RepoDBPostgres, error) {
	return RepoDBPostgres{}, nil
}

// Регистрация пользователя в БД.
//
// Добавляет запись по указанным данным в БД.
//
// Входящие параметры: models.User (параметры Login - регистронезависимый, уникальный)
func (repoPG *RepoDBPostgres) UserRegister(user models.User) error {
	ctx := context.Background()
	_, err := repoPG.PGXPool.Exec(ctx, `INSERT INTO USERS ( name_user, pass_user, created_at, is_active ) 
		 VALUES ( @P_NAME_USER, @P_PASS_USER, CURRENT_TIMESTAMP, 1)`,
		pgx.NamedArgs{"P_NAME_USER": user.Login, "P_PASS_USER": user.Password},
	)

	if pgErr, ok := err.(*pgconn.PgError); ok {
		if pgErr.Code == pgerrcode.UniqueViolation {
			return pgErr
		}
	}
	return err
}

// Аутентификация пользователя в БД.
//
// Проверяет наличие записи с указанными даннымы в БД.
//
// Входящие параметры: models.User (параметры Login - регистронезависимый, уникальный)
func (repoPG *RepoDBPostgres) UserLogin(user models.User) error {
	ctx := context.Background()
	var count int

	err := repoPG.PGXPool.QueryRow(ctx, `SELECT COUNT(1) FROM USERS WHERE name_user = @P_NAME_USER AND PASS_USER = @P_PASS_USER`,
		pgx.NamedArgs{"P_NAME_USER": user.Login, "P_PASS_USER": user.Password},
	).Scan(&count)

	if err != nil {
		return err
	}

	if count == 0 {
		return &pgconn.PgError{
			Code:    "28000",
			Message: "Wrong Username or Password",
		}
	}
	return err
}

// Загрузка номера заказа
//
// Проверяет наличие заказа с указанными даннымы в БД: текущий пользователь/другой пользователь
// Добавляет заказ, если нет
//
// Входящие параметры: models.OrderShort
func (repoPG *RepoDBPostgres) UserOrderCreate(Order models.OrderShort) error {
	ctx := context.Background()
	var orderUser string

	err := repoPG.PGXPool.QueryRow(ctx, `select coalesce(max(name_user), '-1') name_user from users u join orders o on o.user_id = u.user_id where o.order_num = @P_ORDER_NUM`,
		pgx.NamedArgs{"P_ORDER_NUM": Order.OrderNumber},
	).Scan(&orderUser)

	if err != nil {
		return err
	}

	if orderUser == "-1" {
		_, err = repoPG.PGXPool.Exec(ctx, `INSERT INTO ORDERS ( order_num, user_id, total_amount, status, bonus_sum) 
		select @P_ORDER_NUM, user_id, 0, 'CREATED', 0 from users where name_user = @P_USER_NAME`,
			pgx.NamedArgs{"P_ORDER_NUM": Order.OrderNumber, "P_USER_NAME": Order.User},
		)

	} else {
		if orderUser == Order.User {
			return &pgconn.PgError{
				Code:    "-1",
				Message: "The order has already been created by this user",
			}
		} else {
			return &pgconn.PgError{
				Code:    "-2",
				Message: "The order has already been created by another user",
			}
		}
	}

	return err
}

func (repoPG *RepoDBPostgres) GetUserOrders(User string) ([]models.Order, error) {
	clientAccrual := client.NewAccrualClient("http://localhost:8081")
	order, err := clientAccrual.GetOrder("6549842131")
	fmt.Println("order ", order)
	if err != nil {
		fmt.Println("err  ", err)
	}
	var orders []models.Order
	/*
		ctx := context.Background()

		err = repoPG.PGXPool.QueryRow(ctx, `select * from orders o where o.user_id = @P_USER_ID`,
			pgx.NamedArgs{"P_USER_ID": User},
		).Scan(&orders)

		if err != nil {
			return orders, err
		}
	*/
	return orders, nil
}

func (repoPG *RepoDBPostgres) GetBalanceAccrual(User string) (models.OrderAccrual, error) {
	var OrderAccrual models.OrderAccrual
	return OrderAccrual, nil
}

func (repoPG *RepoDBPostgres) GetUserBalance(User string) (models.Balance, error) {
	var Balnce models.Balance
	return Balnce, nil
}

func (repoPG *RepoDBPostgres) RequestOrderAccrual(balanceWD models.BalanceWithDraw) error {
	return nil
}

// Формирование репозитория DB Postgres
//
// Входные параметры: cft config.ConfigGM - конфигурация сервиса строка подключения
func NewRepository(cfg config.ConfigGM) (GopherMart, error) {
	var repoGM RepoDBPostgres
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DSN)
	if err != nil {
		return &repoGM, err
	}
	repoGM.PGXPool = pool
	err = repoGM.CreateDBScheme(ctx, cfg.MigrationPath)

	return &repoGM, err

}

func (repoPG *RepoDBPostgres) Close() {
	repoPG.PGXPool.Close()
}

func (repoPG *RepoDBPostgres) CreateDBScheme(ctx context.Context, MigrationPath string) error {
	conn, err := repoPG.PGXPool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	migrationsDir := MigrationPath
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return err
	}

	for _, file := range files {

		if strings.HasSuffix(file.Name(), ".up.sql") {
			migrationPath := filepath.Join(migrationsDir, file.Name())
			sqlBytes, err := os.ReadFile(migrationPath)
			if err != nil {
				log.Println("Ошибка миграции: ", err.Error())
				continue
			}

			_, err = repoPG.PGXPool.Exec(ctx, string(sqlBytes))
			if err != nil {
				log.Println("Ошибка миграции: ", err.Error())
				continue
			}
		}
	}

	return err
}
