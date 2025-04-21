package repositories

import (
	"context"
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
	//err = pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM users").Scan(&count)

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

func (repoPG *RepoDBPostgres) UserOrderCreate(order string) error {
	return nil
}

func (repoPG *RepoDBPostgres) GetUserOrders() ([]models.Order, error) {
	var orders []models.Order
	return orders, nil
}

func (repoPG *RepoDBPostgres) GetBalanceAccrual() (models.OrderAccrual, error) {
	var OrderAccrual models.OrderAccrual
	return OrderAccrual, nil
}

func (repoPG *RepoDBPostgres) GetUserBalance() (models.Balance, error) {
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
