package middlewares

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	Login string
}

type AuthConfig struct {
	CookieName string
	SecretKey  string
	TokenExp   time.Duration
}

type AuthSvc struct {
	Config AuthConfig
}

func NewAuthService(cfgAuth AuthConfig) AuthSvc {
	return AuthSvc{Config: cfgAuth}
}

// Создаёт токен и возвращает его в виде строки
//
// Используем уникальное имя пользователя
func (auth *AuthSvc) GenerateToken(login string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(auth.Config.TokenExp)),
		},
		Login: login,
	})

	authToken, err := token.SignedString([]byte(auth.Config.SecretKey))
	if err != nil {
		return authToken, err
	}

	return authToken, err
}

// Проверяет токен на валидность
func (auth *AuthSvc) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("NotValid")
		}
		return []byte(auth.Config.SecretKey), nil
	})
}

/*
func WithAuth(authService AuthSvc) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(authService.Config.TokenExp)),
				},
				Login: au.Login,
			})

			authToken, err := token.SignedString([]byte(authService.Config.SecretKey))
			if err != nil {
				return authToken, err
			}

			// Проверяем валидность токена
			login, err := authService.ValidateToken(token)
			if err != nil {

				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
*/
