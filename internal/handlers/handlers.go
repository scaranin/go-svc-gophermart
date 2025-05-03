package handlers

import (
	"go-svc-gophermart/internal/client"
	"go-svc-gophermart/internal/config"
	"go-svc-gophermart/internal/middlewares"
	"go-svc-gophermart/internal/repositories"
	"log"
	"net/http"
	"runtime"
	"unicode"

	"github.com/go-chi/chi"
)

// Основная структура со списком обработчиков
type URLHandler struct {
	Repo       repositories.GopherMart
	TokenSvc   middlewares.TokenService
	AccrualSvc client.AccrualService
}

// Функция для получения имени текущего метода
func (h URLHandler) GetCurrentMethodName() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return "unknown"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}

	fullName := fn.Name()

	return fullName
}

// Обработка cookie
func (h *URLHandler) cookieProcessing(w http.ResponseWriter, r *http.Request) (string, error) {
	var (
		user string
		err  error
	)
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return user, err
	}
	user, err = h.TokenSvc.GetUserFromCookie(cookie)
	if err != nil {
		log.Println(err)
	}

	cookieW, err := h.TokenSvc.GenerateCookie(user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
	}
	http.SetCookie(w, cookieW)

	return user, err
}

// Проверка номера договора алгоритмом Луна
func validateLuhn(orderNumber string) bool {
	for _, r := range orderNumber {
		if !unicode.IsDigit(r) {
			return false
		}
	}

	sum := 0
	for i, rune := range orderNumber {
		digit := int(rune)
		if (len(orderNumber)-i)%2 == 0 {
			digit *= 2

			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}

	return sum%10 == 0
}

// Для тестов
func InitHandlerTest() (URLHandler, error) {
	cfg, err := config.NewConfig()

	if err != nil {
		log.Fatal(err)
	}
	var h URLHandler

	h.Repo, err = repositories.NewRepository(cfg)
	if err != nil {
		log.Fatal(err)
	}

	tokenSvc := middlewares.NewTokenService(cfg.SecretKey)
	h.TokenSvc = &tokenSvc
	accrualSvc := client.NewAccrualClient(cfg.ASAddress)
	accrualSvc.Repo = h.Repo
	h.AccrualSvc = accrualSvc
	accrualSvc.RunTickerWithContext()

	authConfig := middlewares.NewAuthConfig(cfg.SecretKey)

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Up!")

	mux := chi.NewRouter()
	authService := middlewares.NewAuthService(authConfig.SecretKey)

	authMiddleware := middlewares.WithAuth(&authService)

	mux.Use(middlewares.WithLogging, authMiddleware)

	mux.Route("/", func(mux chi.Router) {
		mux.Post("/api/user/register", h.PostUserRegister)
		mux.Post("/api/user/login", h.PostUserLogin)
		mux.Post("/api/user/orders", h.PostUserOrders)
		mux.Post("/api/user/balance/withdraw", h.RequestWithdraw)

		mux.Get("/api/user/orders", h.GetUserOrders)
		mux.Get("/api/user/balance", h.GetUserBalance)
		mux.Get("/api/user/withdrawals", h.GetWithdrawals)
	})

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	go func() {
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	return h, err
}
