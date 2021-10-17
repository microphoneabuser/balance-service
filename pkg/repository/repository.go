package repository

import (
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/microphoneabuser/balance-service/models"
)

type Account interface {
	GetBalance(id int) (int, error)
	AddAccount(id int) error
	Accrual(input models.AccountInput) error
	Debiting(input models.AccountInput) error
	// Update(input models.UpdateAccountInput) error
}

type Transaction interface {
	GetTransaction(transactionId int) (models.Transaction, error)
	MakeTransaction(input models.TransactionInput) error
	GetAccountTransactions(accountId int, outputParams models.OutputParams) ([]models.Transaction, error)
}

type CurrencyAPI interface {
	GetCurrency(code string) (float64, error)
}

type Repository struct {
	Account
	Transaction
	CurrencyAPI
}

func NewRepository(db *sqlx.DB, redisClient *redis.Client) *Repository {
	return &Repository{
		Account:     NewAccountPostgres(db),
		Transaction: NewTransactionPostgres(db),
		CurrencyAPI: NewCurrencyAPIRedis(redisClient),
	}
}
