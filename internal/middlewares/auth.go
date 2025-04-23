package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type TokenService interface {
	GenerateToken(User string) (string, error)
}

type ValidateService interface {
	ValidateToken(tokenString string) (string, error)
}

type AuthService interface {
	TokenService
	ValidateService
}

type Claims struct {
	jwt.RegisteredClaims
	User string
}

type AuthConfig struct {
	CookieName string
	SecretKey  string
	TokenExp   time.Duration
}

func NewAuthConfig(secretKey string) AuthConfig {
	return AuthConfig{
		CookieName: "auth_token",
		SecretKey:  secretKey,
		TokenExp:   time.Hour,
	}
}

type AuthSvc struct {
	Config AuthConfig
}

func NewAuthService(SecretKey string) AuthSvc {
	return AuthSvc{Config: AuthConfig{CookieName: "auth_token",
		SecretKey: SecretKey,
		TokenExp:  time.Hour}}
}

// Создаёт токен и возвращает его в виде строки
//
// Используем уникальное имя пользователя
func (auth *AuthSvc) GenerateToken(User string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(auth.Config.TokenExp)),
		},
		User: User,
	})

	authToken, err := token.SignedString([]byte(auth.Config.SecretKey))
	log.Println(authToken)
	if err != nil {
		return authToken, err
	}

	return authToken, err
}

// Проверяет токен на валидность
func (auth *AuthSvc) ValidateToken(tokenString string) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(auth.Config.TokenExp)),
		},
	}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header[auth.Config.CookieName])
		}
		return auth.Config.SecretKey, nil
	})

	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	return claims.Subject, nil
}

func WithAuth(authService AuthService) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Получаем токен из cookies
			cookie, err := r.Cookie("auth_token")
			if err != nil {
				if err == http.ErrNoCookie {
					next.ServeHTTP(w, r)
					return
				}
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			fmt.Println("cookie", cookie)

			// Проверяем токен с помощью сервиса авторизации
			User, err := authService.ValidateToken(cookie.Value)
			log.Println(User)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Добавляем userID в контекст запроса
			next.ServeHTTP(w, r)
		})
	}
}
