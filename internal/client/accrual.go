package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-svc-gophermart/internal/models"
	"go-svc-gophermart/internal/repositories"
	"log"
	"net/http"
	"time"
)

type AccrualService interface {
	GetOrder(orderNum string) (models.OrderAccrual, error)
}

// Клиент сервиса бонусного счета
type AccrualClient struct {
	baseURL string
	client  *http.Client
	Repo    repositories.GopherMart
}

// Инициализация клиента сервиса бонусного счета
func NewAccrualClient(baseURL string) *AccrualClient {
	return &AccrualClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

// Получение информации о расчёте начислений баллов лояльности по заказу
func (accrual *AccrualClient) GetOrder(orderNum string) (models.OrderAccrual, error) {

	var order models.OrderAccrual

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/api/orders/%s", accrual.baseURL, orderNum), nil)
	if err != nil {
		log.Println(err)
		return order, err
	}

	resp, err := accrual.client.Do(req)
	if err != nil {
		log.Println(err)
		return order, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return order, errors.New(resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&order); err != nil {
		return order, err
	}

	return order, nil
}

// Запуск процесса получения статуса бонусов по заказу
func (accrual *AccrualClient) RunTickerWithContext() error {
	ctx := context.Background()
	go func() {
		ticker := time.NewTicker(time.Second * 5)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("RunTickerWithContext stopped!")
				return
			case <-ticker.C:
				OrderAccrualList, err := accrual.Repo.GetAccrualsFull()
				if err != nil {
					log.Println(err)
				}
				var OrderList []models.OrderAccrual
				for _, OrderAccrualItem := range OrderAccrualList {
					Order, err := accrual.GetOrder(OrderAccrualItem.Order)
					if err != nil {
						log.Println(err)
					}
					OrderList = append(OrderList, Order)
				}
				accrual.Repo.UpdateOrderList(OrderList)
			}
		}
	}()
	return nil
}
