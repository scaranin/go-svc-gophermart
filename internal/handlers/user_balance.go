package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-svc-gophermart/internal/models"
	"log"
	"net/http"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

// Получение текущего баланса счёта баллов лояльности пользователя
func (h *URLHandler) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth_token")
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user, err := h.TokenSvc.GetUserFromCookie(cookie)
	if err != nil {
		log.Println(err)
	}

	Balance, err := h.Repo.GetUserBalance(user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	cookieW, err := h.TokenSvc.GenerateCookie(user)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
	}

	http.SetCookie(w, cookieW)

	BalanceJSON, err := json.Marshal(Balance)
	if err != nil {
		log.Println(err)

		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(BalanceJSON)
}

// Запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа
func (h *URLHandler) RequestWithdraw(w http.ResponseWriter, r *http.Request) {
	var (
		buf      bytes.Buffer
		header   int
		pgErr    *pgconn.PgError
		ok       bool
		Withdraw models.RequestWithDraw
	)
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user, err := h.TokenSvc.GetUserFromCookie(cookie)
	if err != nil {
		log.Println(err)
	}

	defer r.Body.Close()
	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &Withdraw); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	/*
		if validateLuhn(Withdraw.Order) {
			header = http.StatusUnprocessableEntity
			log.Print("Bad order number")
		}
	*/

	Balance, err := h.Repo.GetUserBalance(user)
	Balance.Current = Balance.Current - Balance.WithDrawn
	if err != nil {
		header = http.StatusInternalServerError
		log.Println(err)
		w.WriteHeader(header)
		return
	}

	if Balance.Current-Withdraw.Sum < 0 {
		header = http.StatusPaymentRequired
		log.Print(pgErr)
		w.WriteHeader(header)
		return
	} else {
		fmt.Println("Withdraw ", Withdraw)
		pgErr, ok = h.Repo.RequestOrderAccrual(user, Withdraw).(*pgconn.PgError)
	}

	if ok {
		switch pgErr.Code {
		case pgerrcode.SuccessfulCompletion:
			header = http.StatusOK
			log.Print(pgErr)
		default:
			header = http.StatusInternalServerError
			log.Print(pgErr)
		}
	}

	cookieW, err := h.TokenSvc.GenerateCookie(user)
	if err != nil {
		header = http.StatusUnauthorized
		log.Print(err)
		w.WriteHeader(header)
		return
	}

	http.SetCookie(w, cookieW)
	if header == 0 {
		header = http.StatusAccepted
	}
	w.WriteHeader(header)
}

// Получение информации о выводе средств с накопительного счёта пользователем
func (h *URLHandler) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth_token")
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user, err := h.TokenSvc.GetUserFromCookie(cookie)
	if err != nil {
		log.Println(err)
	}

	// Получаем заказы с актуальным статусом
	WithdrawList, err := h.Repo.GetUserWithdrawAll(user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	cookieW, err := h.TokenSvc.GenerateCookie(user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
	}
	http.SetCookie(w, cookieW)

	WithdrawListJSON, err := json.Marshal(WithdrawList)
	if err != nil {
		log.Println(err)

		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(WithdrawListJSON)
}
