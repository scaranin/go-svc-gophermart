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

// # Регистрация пользователя в БД.
//
// Добавляет запись по указанным данным в БД.
//
// Входные параметры: models.User (параметры Login - регистронезависимый, уникальный)
func (repoPG *RepoDBPostgres) UserRegister(user models.User) error {
	ctx := context.Background()

	SQLUserInsert := `insert 
	                    into USERS ( name_user, pass_user, created_at, is_active ) 
		              values ( @P_NAME_USER, @P_PASS_USER, CURRENT_TIMESTAMP, 1)`

	_, err := repoPG.PGXPool.Exec(ctx, SQLUserInsert,
		pgx.NamedArgs{"P_NAME_USER": user.Login, "P_PASS_USER": user.Password},
	)

	if pgErr, ok := err.(*pgconn.PgError); ok {
		if pgErr.Code == pgerrcode.UniqueViolation {
			return pgErr
		}
	}
	return err
}

// # Аутентификация пользователя в БД.
//
// Проверяет наличие записи с указанными даннымы в БД.
//
// Входные параметры: models.User (параметры Login - регистронезависимый, уникальный)
func (repoPG *RepoDBPostgres) UserLogin(user models.User) error {
	ctx := context.Background()
	var count int

	SQLExists := `select count(1) 
	                form USERS 
				   where name_user = @P_NAME_USER 
				     and pass_user = @P_PASS_USER`

	err := repoPG.PGXPool.QueryRow(ctx, SQLExists,
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

// # Загрузка номера заказа
//
// Проверяет наличие заказа с указанными даннымы в БД: текущий пользователь/другой пользователь
// Добавляет заказ, если нет
//
// Входные параметры: models.OrderShort
func (repoPG *RepoDBPostgres) UserOrderCreate(Order models.OrderShort) error {
	ctx := context.Background()
	var orderUser string

	SQLExists := `select coalesce(max(U.name_user), '-1') name_user 
	                from USERS  U 
				    join ORDERS O on O.user_id = U.user_id 
				   where O.order_num = @P_ORDER_NUM`

	err := repoPG.PGXPool.QueryRow(ctx, SQLExists,
		pgx.NamedArgs{"P_ORDER_NUM": Order.OrderNumber},
	).Scan(&orderUser)

	if err != nil {
		return err
	}

	SQLInsertOrder := `insert 
	                     into ORDERS ( order_num, user_id, total_amount, status, accrual) 
		               select @P_ORDER_NUM
					        , user_id
							, 0
							, 'NEW'
							, 0 
						 from users 
						where name_user = @P_USER_NAME`

	if orderUser == "-1" {
		_, err = repoPG.PGXPool.Exec(ctx, SQLInsertOrder,
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

// # Получение всех заказов не в конечных статусах
//
// Отбирает заказы для отправки запроса в сервис Accrual
func (repoPG *RepoDBPostgres) GetAccrualsFull() ([]models.OrderAccrual, error) {
	ctx := context.Background()
	var (
		err              error
		orderAccrualList []models.OrderAccrual
		orderAccrealItem models.OrderAccrual
	)
	SQLGetOrderList :=
		`select O.order_num
           from ORDERS             O
           join USERS              U on U.user_id = O.user_id
           join STATUS_ORDER_INFO SI on SI.status = O.status
          where CAST(NOW() AS DATE) between SI.date_start
                                        and SI.date_end
            and SI.is_active = 1
            and SI.is_final  = 0`

	rows, err := repoPG.PGXPool.Query(ctx, SQLGetOrderList)
	if err != nil {
		log.Println(err)
		return orderAccrualList, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&orderAccrealItem.Order)
		if err != nil {
			log.Println(err)
		}
		orderAccrualList = append(orderAccrualList, orderAccrealItem)

	}

	return orderAccrualList, err
}

// # Получение списка заказов пользователя не в конечных статусах
//
//	Отбирает заказы пользователя в статусе отличном от конечного
//
// Входные параметры: User
func (repoPG *RepoDBPostgres) GetAccruals(User string) ([]models.OrderAccrual, error) {
	ctx := context.Background()
	var (
		err              error
		orderAccrualList []models.OrderAccrual
		orderAccrealItem models.OrderAccrual
	)

	SQLGetOrderList :=
		`select O.order_num
           from ORDERS             O
           join USERS              U on U.user_id = O.user_id
           join STATUS_ORDER_INFO SI on SI.status = O.status
          where CAST(NOW() AS DATE) between SI.date_start
                                        and SI.date_end
            and SI.is_active = 1
            and SI.is_final  = 0
            and U.name_user  = @P_USER_NAME`

	rows, err := repoPG.PGXPool.Query(ctx, SQLGetOrderList, pgx.NamedArgs{"P_USER_NAME": User})
	if err != nil {
		log.Println(err)
		return orderAccrualList, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&orderAccrealItem.Order)
		if err != nil {
			log.Println(err)
		}
		orderAccrualList = append(orderAccrualList, orderAccrealItem)

	}

	return orderAccrualList, err

}

func (repoPG *RepoDBPostgres) UpdateOrderList(OrderAccrualArr []models.OrderAccrual) error {
	ctx := context.Background()
	tx, err := repoPG.PGXPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	SQLUpdateOrder := `update ORDERS 
	                      set status     = $1
						    , accrual    = $2
							, updated_at = current_timestamp 
						where order_num = $3`

	_, err = tx.Prepare(ctx, "UploadAccrual", SQLUpdateOrder)
	if err != nil {
		return err
	}

	for _, OrderAccrual := range OrderAccrualArr {
		_, err := tx.Exec(ctx, "UploadAccrual", OrderAccrual.Status, OrderAccrual.Accrual, OrderAccrual.Order)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// # Получение списка заказов пользователя
//
//	Возвращает список заказов: текущий пользователь/другой пользователь
//
// Входные параметры: User
// Выходные параметры: []models.Order
func (repoPG *RepoDBPostgres) GetUserOrders(User string) ([]models.Order, error) {
	var orderList []models.Order
	var order models.Order
	ctx := context.Background()

	sqlGetOrderList := `select O.order_num
	                         , O.status
							 , O.accrual
							 , O.updated_at 
						  from ORDERS O
						  join USERS  U on U.user_id = O.user_id 
						 where U.name_user = @P_USER_NAME`

	var err error
	rows, err := repoPG.PGXPool.Query(ctx, sqlGetOrderList, pgx.NamedArgs{"P_USER_NAME": User})
	if err != nil {
		return orderList, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.DtUploaded)
		if err != nil {
			log.Println(err)
		}

		orderList = append(orderList, order)
	}

	if err != nil {
		return orderList, err
	}

	return orderList, err
}

// # Получение списка заказов пользователя. Внешний сервис
//
//	Возвращает список заказов пользователя влючая статус и начисленные баллы
//
// Входные параметры: User
func (repoPG *RepoDBPostgres) GetBalanceAccrual(User string) ([]models.OrderAccrual, error) {
	var OrderAccrualList []models.OrderAccrual
	return OrderAccrualList, nil
}

func (repoPG *RepoDBPostgres) GetUserBalance(User string) (models.Balance, error) {
	ctx := context.Background()
	var Balnce models.Balance

	SQLSelectBalance := `select ( select sum(accrual) from ORDERS O where O.user_id = U.user_id )   accrual
	                          , ( select sum(amount) from WITHDRAWS W where W.user_id = U.user_id ) withdraw
	                       from USERS U 
				          where U.user_name = @P_USER_NAME`

	err := repoPG.PGXPool.QueryRow(ctx, SQLSelectBalance,
		pgx.NamedArgs{"P_USER_NAME": User},
	).Scan(&Balnce.Current, &Balnce.WithDrawn)

	return Balnce, err

}

// # Запрос на списание средств
//
//	Создаем запрос на списание бонусных средств в счет заказа
//
// Входные параметры: models.RequestWithDraw
func (repoPG *RepoDBPostgres) RequestOrderAccrual(User string, WithDraw models.RequestWithDraw) error {
	ctx := context.Background()

	SQLInsertRequest := `insert 
	                       into WITHDRAWS ( order_num, user_id, amount, status, updated_at ) 
		                 select @P_ORDER_NUM, user_id, 0, 'REGISTERED', current_timestamp from users where name_user = @P_USER_NAME`

	_, err := repoPG.PGXPool.Exec(ctx, SQLInsertRequest,
		pgx.NamedArgs{"P_ORDER_NUM": WithDraw.Order, "P_SUM_REQUEST": WithDraw.Sum, "P_USER_NAME": User},
	)

	if err != nil {
		return err
	}

	return err

}

// GetUserWithdrawAll implements GopherMart.
func (repoPG *RepoDBPostgres) GetUserWithdrawAll(User string) ([]models.WithDrawalsList, error) {
	var (
		withdrawList []models.WithDrawalsList
		withdraw     models.WithDrawalsList
		err          error
	)
	ctx := context.Background()

	sqlGetWithdrawList := `select O.order_num
	                            , O.amount
						        , O.updated_at 
						     from WITHDRAWS O
						     join USERS     U on U.user_id = O.user_id 
						    where U.name_user = @P_USER_NAME`

	rows, err := repoPG.PGXPool.Query(ctx, sqlGetWithdrawList, pgx.NamedArgs{"P_USER_NAME": User})
	if err != nil {
		return withdrawList, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&withdraw.Order, &withdraw.Sum, &withdraw.DtProcessed)
		if err != nil {
			log.Println(err)
		}

		withdrawList = append(withdrawList, withdraw)
	}

	return withdrawList, err
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

			pgErr, ok := err.(*pgconn.PgError)

			if ok && pgErr.Code != pgerrcode.DuplicateTable {
				log.Println("Ошибка миграции: ", err.Error())
				continue
			}
		}
	}

	return err
}
