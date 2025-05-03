package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-svc-gophermart/internal/models"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURLHandler_GetUserOrders(t *testing.T) {
	type want struct {
		statusCode int
		request    string
		login      string
		password   string
	}

	login, err := GenerateRandomString(5)
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
			name: "GetUserOrders test #1",
			want: want{
				statusCode: http.StatusNoContent,
				request:    "http://localhost:8080/api/user/orders",
				login:      login,
				password:   password,
			},
		},
	}

	h, err := InitHandlerTest()

	if err != nil {
		log.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var user models.User = models.User{Login: tt.want.login, Password: tt.want.password}

			jsonUser, err := json.Marshal(user)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(string(jsonUser))

			resp, err := http.Post("http://localhost:8080/api/user/register", "application/json", bytes.NewBuffer(jsonUser))
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
				assert.Equal(t, resp.StatusCode, http.StatusOK)
			}

			req, err := http.NewRequest("GET", tt.want.request, bytes.NewBuffer([]byte(user.Login)))
			if err != nil {
				fmt.Println(err)
				return
			}

			req.Header.Set("Content-Type", "application/json")

			cookieW, err := h.TokenSvc.GenerateCookie(tt.want.login)
			if err != nil {
				log.Fatal(err)
			}
			req.AddCookie(cookieW)

			client := &http.Client{}

			respOrder, err := client.Do(req)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer resp.Body.Close()

			bodyOrder, err := io.ReadAll(respOrder.Body)
			if err != nil {
				fmt.Println(err)
				return
			}

			if bodyOrder != nil {
				assert.Equal(t, tt.want.statusCode, respOrder.StatusCode)
			}
		})
	}
}
