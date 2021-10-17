package service

import (
	"github.com/microphoneabuser/balance-service/models"
	"github.com/microphoneabuser/balance-service/pkg/repository"
)

type Account interface {
	GetBalance(id int) (int, error)
	Accrual(input models.AccountInput) error
	Debiting(input models.AccountInput) error
}

type Transaction interface {
	GetTransaction(transactionId int) (models.Transaction, error)
	MakeTransaction(input models.TransactionInput) error
	GetAccountTransactions(accountId int, outputParams models.OutputParams) ([]models.TransactionForOutput, error)
}

type Currency interface {
	Convert(amount int, code string) (int, error)
}

type Service struct {
	Account
	Transaction
	Currency
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Account:     NewAccountService(repos.Account),
		Transaction: NewTransactionService(repos.Transaction),
		Currency:    NewCurrencyService(repos.CurrencyAPI),
	}
}
