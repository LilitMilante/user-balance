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

func (s *Store) InsertUserBalance(b entity.Balance) error {
	_, err := s.db.Exec("INSERT INTO users_balances (user_id, amount) VALUES ($1, $2)", b.UserID, b.Amount)
	if err != nil {
		return err
	}

	return nil
}

//
//func (s *Store) UpdateUserBalance(id int64, sum int64) (int64, error) {
//	var b entity.Balance
//
//	err := s.db.QueryRow("UPDATE users_balances SET amount = amount + $1 WHERE user_id = $2 RETURNING amount", sum, id).Scan(&b.Amount)
//	if err != nil {
//		return b.Amount, err
//	}
//
//	b.UserID = id
//
//	return b.Amount, nil
//}

func (s *Store) TxUpdateUsersBalances(t entity.Transfer) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	var curSum int64

	err = tx.QueryRow("SELECT amount FROM users_balances WHERE user_id = $1", t.IdTake).Scan(&curSum)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.ErrNotFound
		}

		return err
	}

	if (curSum - t.Amount) < 0 {
		return domain.ErrEnoughMoney
	}

	_, err = tx.Exec("UPDATE users_balances SET amount = amount - $1 WHERE user_id = $2", t.Amount, t.IdGive)
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO users_transactions (user_id, amount, type_op, description) VALUES ($1, $2, $3, $4)", t.IdGive, t.Amount, entity.Minus, t.DescriptionSender)
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE users_balances SET amount = amount + $1 WHERE user_id = $2", t.Amount, t.IdTake)
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO users_transactions (user_id, amount, type_op, description) VALUES ($1, $2, $3, $4)", t.IdTake, t.Amount, entity.Plus, t.DescriptionRecipient)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) SelectUserTransactions(uID int64) ([]entity.Balance, error) {
	var b []entity.Balance

	rows, err := s.db.Query("SELECT user_id, amount, type_op, description FROM users_transactions WHERE user_id = $1", uID)
	if err != nil {

		return nil, err
	}

	defer rows.Close()

	if !rows.NextResultSet() {
		return nil, domain.ErrNotFound
	}

	for rows.Next() {
		var t entity.Balance

		err := rows.Scan(&t.UserID, &t.Amount, &t.TypeOp, &t.Description)
		if err != nil {
			return nil, err
		}

		b = append(b, t)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *Store) InsertUserTransactions(b entity.Balance) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = s.db.Exec("INSERT INTO users_balances (user_id, amount) VALUES ($1, $2)", b.UserID, b.Amount)
	if err != nil {
		return err
	}

	_, err = s.db.Exec("INSERT INTO users_transactions (user_id, amount, type_op, description) VALUES ($1, $2, $3, $4)", b.UserID, b.Amount, b.TypeOp, b.Description)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdateUserTransactions(b entity.Balance) (int64, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return b.Amount, err
	}

	defer tx.Rollback()

	err = s.db.QueryRow("UPDATE users_balances SET amount = amount + $1 WHERE user_id = $2 RETURNING amount", b.Amount, b.UserID).Scan(&b.Amount)
	if err != nil {
		return b.Amount, err
	}

	_, err = s.db.Exec("INSERT INTO users_transactions (user_id, amount, type_op, description) VALUES ($1, $2, $3, $4)", b.UserID, b.Amount, b.TypeOp, b.Description)
	if err != nil {
		return b.Amount, err
	}

	err = tx.Commit()
	if err != nil {
		return b.Amount, err
	}

	return b.Amount, nil
}
