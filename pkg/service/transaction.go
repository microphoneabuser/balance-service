package service

import (
	"fmt"

	"github.com/microphoneabuser/balance-service/models"
	"github.com/microphoneabuser/balance-service/pkg/repository"
)

type TransactionService struct {
	repo repository.Transaction
}

func NewTransactionService(repo repository.Transaction) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) GetTransaction(transactionId int) (models.Transaction, error) {
	return s.repo.GetTransaction(transactionId)
}
func (s *TransactionService) MakeTransaction(input models.TransactionInput) error {
	if input.Amount <= 0 {
		return fmt.Errorf("The amount must be greater than zero")
	}

	err := s.repo.MakeTransaction(input)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return fmt.Errorf("Account not found")
	}

	return err
}
func (s *TransactionService) GetAccountTransactions(accountId int, outputParams models.OutputParams) ([]models.TransactionForOutput, error) {
	output := []models.TransactionForOutput{}
	trs, err := s.repo.GetAccountTransactions(accountId, outputParams)
	if err != nil {
		return output, err
	}

	for _, transaction := range trs {
		switch {
		case transaction.RecipientId == 0:
			output = append(output, models.TransactionForOutput{
				Id:          transaction.Id,
				TimeStamp:   transaction.TimeStamp,
				Type:        "Списание",
				Description: transaction.Description,
				Amount:      toNormal(transaction.Amount),
			})
		case transaction.SenderId == 0:
			output = append(output, models.TransactionForOutput{
				Id:          transaction.Id,
				TimeStamp:   transaction.TimeStamp,
				Type:        "Зачисление",
				Description: transaction.Description,
				Amount:      toNormal(transaction.Amount),
			})
		case transaction.RecipientId == accountId:
			output = append(output, models.TransactionForOutput{
				Id:          transaction.Id,
				TimeStamp:   transaction.TimeStamp,
				Type:        "Входящий перевод",
				SenderId:    transaction.SenderId,
				Description: transaction.Description,
				Amount:      toNormal(transaction.Amount),
			})
		case transaction.SenderId == accountId:
			output = append(output, models.TransactionForOutput{
				Id:          transaction.Id,
				TimeStamp:   transaction.TimeStamp,
				Type:        "Исходящий перевод",
				RecipientId: transaction.RecipientId,
				Description: transaction.Description,
				Amount:      toNormal(transaction.Amount),
			})
		}
	}

	return output, nil
}

func toNormal(amount int) float64 {
	return float64(amount) / 100
}
