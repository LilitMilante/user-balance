package store

import (
	"database/sql"
	"errors"
	"user-balance/domain"
	"user-balance/domain/entity"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	s := Store{
		db: db,
	}

	return &s
}

func (s *Store) InsertUserTransactions(t entity.Transaction) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO users_balances (user_id, amount) VALUES ($1, $2)", t.UserID, t.Amount)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT INTO users_transactions (user_id, amount, type_op, transfer_id, description, created_at) 
								VALUES ($1, $2, $3, $4, $5, $6)`, t.UserID, t.Amount, t.Event, t.TransferID, t.Message, t.CreatedAt)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) InsertUserBalance(b entity.Balance) error {
	_, err := s.db.Exec("INSERT INTO users_balances (user_id, amount) VALUES ($1, $2)", b.UserID, b.Amount)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) SelectUserTransactions(uID int64) ([]entity.Transaction, error) {
	var ts []entity.Transaction

	rows, err := s.db.Query("SELECT id, user_id, amount, type_op, transfer_id, description, created_at FROM users_transactions WHERE user_id = $1", uID)
	if err != nil {

		return nil, err
	}

	defer rows.Close()

	if !rows.NextResultSet() {
		return nil, domain.ErrNotFound
	}

	for rows.Next() {
		var t entity.Transaction

		err := rows.Scan(&t.ID, &t.UserID, &t.Amount, &t.Event, &t.TransferID, &t.Message, &t.CreatedAt)
		if err != nil {
			return nil, err
		}

		ts = append(ts, t)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return ts, nil
}

func (s *Store) SelectUserBalanceByID(id int64) (entity.Balance, error) {
	var b entity.Balance

	err := s.db.QueryRow("SELECT user_id, amount FROM users_balances WHERE user_id = $1", id).Scan(&b.UserID, &b.Amount)

	if errors.Is(err, sql.ErrNoRows) {
		return b, domain.ErrNotFound
	}

	if err != nil {
		return b, err
	}

	return b, nil
}

func (s *Store) UpdateUserTransactions(t entity.Transaction) (int64, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return t.Amount, err
	}

	defer tx.Rollback()

	err = tx.QueryRow("UPDATE users_balances SET amount = amount + $1 WHERE user_id = $2 RETURNING amount", t.Amount, t.UserID).Scan(&t.Amount)
	if err != nil {
		return t.Amount, err
	}

	_, err = tx.Exec(`INSERT INTO users_transactions (user_id, amount, type_op, transfer_id, description, created_at) 
								VALUES ($1, $2, $3, $4, $5, $6)`, t.UserID, t.Amount, t.Event, t.TransferID, t.Message, t.CreatedAt)
	if err != nil {
		return t.Amount, err
	}

	err = tx.Commit()
	if err != nil {
		return t.Amount, err
	}

	return t.Amount, nil
}

func (s *Store) TxUpdateUsersBalances(t entity.Transaction) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	var curSum int64

	err = tx.QueryRow("SELECT amount FROM users_balances WHERE user_id = $1", t.TransferID).Scan(&curSum)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.ErrNotFound
		}

		return err
	}

	if (curSum - t.Amount) < 0 {
		return domain.ErrEnoughMoney
	}

	_, err = tx.Exec("UPDATE users_balances SET amount = amount - $1 WHERE user_id = $2", t.Amount, t.TransferID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT INTO users_transactions (user_id, amount, type_op, transfer_id, description, created_at) 
								VALUES ($1, $2, $3, $4, $5, $6)`, t.TransferID, t.Amount, t.Event, t.UserID, t.Message, t.CreatedAt)
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE users_balances SET amount = amount + $1 WHERE user_id = $2", t.Amount, t.UserID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT INTO users_transactions (user_id, amount, type_op, transfer_id, description, created_at) 
								VALUES ($1, $2, $3, $4, $5, $6)`, t.UserID, t.Amount, t.Event, t.TransferID, t.Message, t.CreatedAt)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
