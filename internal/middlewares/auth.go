package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type TokenService interface {
	GenerateCookie(User string) (*http.Cookie, error)
	GetUserFromCookie(cookie *http.Cookie) (string, error)
}

type TokenSvc struct {
	Config AuthConfig
}

func NewTokenService(SecretKey string) TokenSvc {
	return TokenSvc{Config: AuthConfig{CookieName: "auth_token",
		SecretKey: SecretKey,
		TokenExp:  time.Hour}}
}

type AuthService interface {
	ValidateToken(tokenString string) (string, error)
}

type Claims struct {
	jwt.RegisteredClaims
	User string `json:"user,omitempty"`
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

// Создаёт куку
//
// Используем уникальное имя пользователя
func (tokenSvc *TokenSvc) GenerateCookie(User string) (*http.Cookie, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenSvc.Config.TokenExp)),
		},
		User: User,
	})

	authToken, err := token.SignedString([]byte(tokenSvc.Config.SecretKey))

	if err != nil {
		return nil, err
	}

	cookieW := &http.Cookie{
		Name:     tokenSvc.Config.CookieName,
		Value:    authToken,
		Expires:  time.Now().Add(tokenSvc.Config.TokenExp),
		HttpOnly: true,
		Path:     "/",
	}

	return cookieW, err
}

// Получение User из куки
func (tokenSvc *TokenSvc) GetUserFromCookie(cookie *http.Cookie) (string, error) {
	var (
		user string
		err  error
	)

	if cookie != nil {
		token := cookie.Value
		claims := &Claims{}

		_, err = jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(tokenSvc.Config.SecretKey), nil
		})

		if err != nil {
			return user, err
		}

		user = claims.User
	}
	return user, err
}

// Проверяет токен на валидность
func (auth *AuthSvc) ValidateToken(tokenString string) (string, error) {
	log.Println("tokenString", tokenString)
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(auth.Config.SecretKey), nil
	})

	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	return claims.Subject, nil
}

func WithAuth(authService AuthService) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			cookie, err := r.Cookie("auth_token")
			if err != nil {
				if err == http.ErrNoCookie {
					next.ServeHTTP(w, r)
					return
				}
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			_, err = authService.ValidateToken(cookie.Value)
			if err != nil {
				http.Error(w, "Invalid token"+err.Error(), http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
