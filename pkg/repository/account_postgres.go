package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/microphoneabuser/balance-service/models"
)

type AccountPostgres struct {
	db *sqlx.DB
}

func NewAccountPostgres(db *sqlx.DB) *AccountPostgres {
	return &AccountPostgres{db: db}
}

func (r *AccountPostgres) GetBalance(id int) (int, error) {
	var balance int

	query := fmt.Sprintf("SELECT balance FROM %s WHERE id = $1", accountsTable)

	if err := r.db.Get(&balance, query, id); err != nil {
		return balance, err
	}

	return balance, nil
}

func (r *AccountPostgres) AddAccount(id int) error {
	query := fmt.Sprintf("INSERT INTO %s (id, balance) VALUES ($1, 0)", accountsTable)

	if _, err := r.db.Query(query, id); err != nil {
		return err
	}

	return nil
}

func (r *AccountPostgres) Accrual(input models.AccountInput) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	updateQuery := fmt.Sprintf("UPDATE %s SET balance = balance + $2 WHERE id = $1", accountsTable)

	if _, err := tx.Exec(updateQuery, input.Id, input.Amount); err != nil {
		tx.Rollback()
		return err
	}

	insertQuery := fmt.Sprintf("INSERT INTO %s (recipient_id, amount, description) VALUES ($1, $2, $3)", transactionsTable)

	_, err = tx.Exec(insertQuery, input.Id, input.Amount, input.Description)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *AccountPostgres) Debiting(input models.AccountInput) error {
	currBalance, err := r.GetBalance(input.Id)
	if err != nil {
		return err
	}
	//проверка хватает ли денег на счету
	if currBalance-input.Amount < 0 {
		return fmt.Errorf("Insufficient funds on account %d (id=%d)", input.Id, input.Id)
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	updateQuery := fmt.Sprintf("UPDATE %s SET balance = balance - $2 WHERE id = $1 RETURNING id, balance", accountsTable)

	if _, err := tx.Exec(updateQuery, input.Id, input.Amount); err != nil {
		tx.Rollback()
		return err
	}

	insertQuery := fmt.Sprintf("INSERT INTO %s (sender_id, amount, description) VALUES ($1, $2, $3)", transactionsTable)

	_, err = tx.Exec(insertQuery, input.Id, input.Amount, input.Description)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// func (r *AccountPostgres) Update(input models.UpdateAccountInput) error {
// 	tx, err := r.db.Begin()
// 	if err != nil {
// 		return err
// 	}

// 	updateQuery := fmt.Sprintf("UPDATE %s SET balance = $2 WHERE id = $1 RETURNING id, balance", accountsTable)

// 	if _, err := tx.Query(updateQuery, input.Id, input.Balance); err != nil {
// 		tx.Rollback()
// 		return err
// 	}

// 	var insertQuery string
// 	if input.IsSender {
// 		insertQuery = fmt.Sprintf("INSERT INTO %s (sender_id, amount, description) VALUES ($1, $2, $3)", transactionsTable)
// 	} else {
// 		insertQuery = fmt.Sprintf("INSERT INTO %s (recipient_id, amount, description) VALUES ($1, $2, $3)", transactionsTable)
// 	}

// 	_, err = tx.Exec(insertQuery, input.Id, input.Amount, input.Description)
// 	if err != nil {
// 		tx.Rollback()
// 		return err
// 	}

// 	return nil
// }
