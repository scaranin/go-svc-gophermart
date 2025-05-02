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
	w.Header().Set("Content-Type", "application/json")

	user, err := h.cookieProcessing(w, r)

	if err != nil {
		log.Println(err)
		return
	}

	Balance, err := h.Repo.GetUserBalance(user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	Balance.Current = Balance.Current - Balance.WithDrawn

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
	user, err := h.cookieProcessing(w, r)

	if err != nil {
		log.Println(err)
		return
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

	w.WriteHeader(header)
}

// Получение информации о выводе средств с накопительного счёта пользователем
func (h *URLHandler) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := h.cookieProcessing(w, r)

	if err != nil {
		log.Println(err)
		return
	}

	// Получаем заказы с актуальным статусом
	WithdrawList, err := h.Repo.GetUserWithdrawAll(user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	WithdrawListJSON, err := json.Marshal(WithdrawList)
	if err != nil {
		log.Println(err)

		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(WithdrawListJSON)
}
