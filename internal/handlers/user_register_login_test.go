package handlers

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go-svc-gophermart/internal/models"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func GenerateRandomString(length int) (string, error) {
	byteSlice := make([]byte, length)

	_, err := rand.Read(byteSlice)
	if err != nil {
		return "", err
	}

	randomString := base64.RawStdEncoding.EncodeToString(byteSlice)

	return randomString[:length], nil
}

func TestURLHandler_PostUserRegister(t *testing.T) {
	type want struct {
		statusCode int
		request    string
		login      string
		password   string
	}

	login, err := GenerateRandomString(10)
	if err != nil {
		fmt.Println(err)
	}
	password, err := GenerateRandomString(10)
	if err != nil {
		fmt.Println(err)
	}

	tests := []struct {
		name string
		want want
	}{
		{
			name: "PostUserRegisterAndLogin test #1",
			want: want{
				statusCode: http.StatusOK,
				request:    "http://localhost:8080/api/user/register",
				login:      login,
				password:   password,
			},
		},
		{
			name: "PostUserRegisterAndLogin test #2",
			want: want{
				statusCode: http.StatusNoContent,
				request:    "http://localhost:8080/api/user/login",
				login:      login,
				password:   password,
			},
		},
	}

	if _, err := InitHandlerTest(); err != nil {
		log.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := models.User{Login: tt.want.login, Password: tt.want.password}

			jsonUser, err := json.Marshal(user)
			if err != nil {
				fmt.Println(err)
				return
			}

			resp, err := http.Post(tt.want.request, "application/json", bytes.NewBuffer(jsonUser))
			if err != nil {
				fmt.Println(err)
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
				return
			}
			if body == nil {
				assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			}
		})
	}
}
