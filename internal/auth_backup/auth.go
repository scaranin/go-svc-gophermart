package auth_backup

import (
	"net/http"
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
	Login      string
}

func NewAuthConfig(secretKey string) AuthConfig {
	return AuthConfig{
		CookieName: "auth_token",
		SecretKey:  secretKey,
		TokenExp:   time.Hour,
	}
}

// Создаёт токен и возвращает его в виде строки
//
// Используем уникальное имя пользователя
func (auth *AuthConfig) BuildJWTString() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(auth.TokenExp)),
		},
		Login: auth.Login,
	})

	authToken, err := token.SignedString([]byte(auth.SecretKey))
	if err != nil {
		return authToken, err
	}

	return authToken, err
}

// Формируем ответный coockie. Login записывается в auth.Login
func (auth *AuthConfig) FillUserReturnCookie(incomeCookie *http.Cookie) (*http.Cookie, error) {
	var (
		resAuthToken string
		err          error
	)

	if incomeCookie != nil {
		resAuthToken = incomeCookie.Value
	}

	if len(resAuthToken) == 0 {
		resAuthToken, err = auth.BuildJWTString()
	}
	claims := &Claims{}

	jwt.ParseWithClaims(resAuthToken, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(auth.SecretKey), nil
	})

	if len(claims.Login) == 0 {
		err = http.ErrNoCookie
	} else {
		auth.Login = claims.Login
	}
	cookie := &http.Cookie{
		Name:     auth.CookieName,
		Value:    resAuthToken,
		Expires:  time.Now().Add(auth.TokenExp),
		HttpOnly: true,
		Path:     "/",
	}
	return cookie, err
}

func (auth *AuthConfig) FillUserCookie(Login string) (*http.Cookie, error) {

	resAuthToken, err := auth.BuildJWTString()
	if err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name:     auth.CookieName,
		Value:    resAuthToken,
		Expires:  time.Now().Add(auth.TokenExp),
		HttpOnly: true,
		Path:     "/",
	}
	return cookie, err
}
