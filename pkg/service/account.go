package service

import (
	"fmt"

	"github.com/microphoneabuser/balance-service/models"
	"github.com/microphoneabuser/balance-service/pkg/repository"
)

type AccountService struct {
	repo repository.Account
}

func NewAccountService(repo repository.Account) *AccountService {
	return &AccountService{repo: repo}
}

func (s *AccountService) GetBalance(id int) (int, error) {
	balance, err := s.repo.GetBalance(id)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return 0, fmt.Errorf("There is no account with id=%d", id)
	}
	return balance, err
}
func (s *AccountService) Accrual(input models.AccountInput) error {
	if input.Amount <= 0 {
		return fmt.Errorf("The amount must be greater than zero")
	}

	_, err := s.repo.GetBalance(input.Id)
	if err != nil && err.Error() == "sql: no rows in result set" {
		if err := s.repo.AddAccount(input.Id); err != nil {
			return err
		}
	}

	return s.repo.Accrual(input)
}
func (s *AccountService) Debiting(input models.AccountInput) error {
	if input.Amount <= 0 {
		return fmt.Errorf("The amount must be greater than zero")
	}

	err := s.repo.Debiting(input)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return fmt.Errorf("There is no account with id=%d", input.Id)
	}

	return err
}
