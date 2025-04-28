package client

import (
	"encoding/json"
	"fmt"
	"go-svc-gophermart/internal/models"
	"log"
	"net/http"
)

type AccrualClient struct {
	baseURL string
	client  *http.Client
}

func NewAccrualClient(baseURL string) *AccrualClient {
	return &AccrualClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (accrual *AccrualClient) GetOrder(orderNum string) (*models.OrderAccrual, error) {
	log.Println("GetOrder")
	resp, err := accrual.client.Get(fmt.Sprintf("%s/api/orders/%s", accrual.baseURL, orderNum))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	log.Println("resp.StatusCode ", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HttpStatus: %s", resp.Status)
	}

	var order models.OrderAccrual
	if err := json.NewDecoder(resp.Body).Decode(&order); err != nil {
		return nil, err
	}

	return &order, nil
}
