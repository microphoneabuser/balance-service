package service

import (
	"github.com/microphoneabuser/balance-service/pkg/repository"
)

type CurrencyService struct {
	repo repository.CurrencyAPI
}

func NewCurrencyService(repo repository.CurrencyAPI) *CurrencyService {
	return &CurrencyService{repo: repo}
}

func (c *CurrencyService) Convert(amount int, code string) (int, error) {
	rate, err := c.repo.GetCurrency(code)
	if err != nil {
		return 0, err
	}
	rateInt := int(rate * 1000000)
	return (amount * rateInt) / 1000000, nil
}
