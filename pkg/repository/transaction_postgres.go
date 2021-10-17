package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/microphoneabuser/balance-service/models"
)

type TransactionPostgres struct {
	db *sqlx.DB
}

func NewTransactionPostgres(db *sqlx.DB) *TransactionPostgres {
	return &TransactionPostgres{db: db}
}

func (r *TransactionPostgres) GetTransaction(transactionId int) (models.Transaction, error) {
	var tr models.Transaction

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", transactionsTable)

	if err := r.db.Get(&tr, query, transactionId); err != nil {
		return tr, err
	}

	return tr, nil
}
func (r *TransactionPostgres) MakeTransaction(input models.TransactionInput) error {
	var balance int

	selectQuery := fmt.Sprintf("SELECT balance FROM %s WHERE id = $1", accountsTable)

	if err := r.db.Get(&balance, selectQuery, input.SenderId); err != nil {
		return err
	}

	if balance-input.Amount < 0 {
		return fmt.Errorf("Insufficient funds on account %d (id=%d)", input.SenderId, input.SenderId)
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	updateSenderQuery := fmt.Sprintf("UPDATE %s SET balance = balance - $2 WHERE id = $1", accountsTable)

	if _, err := tx.Exec(updateSenderQuery, input.SenderId, input.Amount); err != nil {
		tx.Rollback()
		return err
	}

	updateRecipientQuery := fmt.Sprintf("UPDATE %s SET balance = balance + $2 WHERE id = $1", accountsTable)

	if _, err := tx.Exec(updateRecipientQuery, input.RecipientId, input.Amount); err != nil {
		tx.Rollback()
		return err
	}

	insertQuery := fmt.Sprintf("INSERT INTO %s (sender_id, recipient_id, amount, description) VALUES ($1, $2, $3, $4)", transactionsTable)

	_, err = tx.Exec(insertQuery, input.SenderId, input.RecipientId, input.Amount, input.Description)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *TransactionPostgres) GetAccountTransactions(accountId int, outputParams models.OutputParams) ([]models.Transaction, error) {
	var trs []models.Transaction

	query := fmt.Sprintf(`SELECT id, COALESCE (sender_id, 0) AS sender_id, COALESCE (recipient_id, 0) AS recipient_id, amount, description, timestamp 
		FROM %s WHERE sender_id = $1 OR recipient_id = $1
		ORDER BY %s %s
		LIMIT $2
		OFFSET $3`, transactionsTable, outputParams.SortCol, outputParams.SortDir)

	if err := r.db.Select(&trs, query, accountId, outputParams.Limit, outputParams.Offset); err != nil {
		return trs, err
	}

	return trs, nil
}
