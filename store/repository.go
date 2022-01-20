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

func (s *Store) SelectUserByID(id int) (entity.Balance, error) {
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

func (s *Store) UpdateUserBalance(id, sum int) (entity.Balance, error) {
	var b entity.Balance

	err := s.db.QueryRow("UPDATE users_balances SET amount = amount + $1 WHERE user_id = $2 RETURNING amount", sum, id).Scan(&b.Amount)
	if err != nil {
		return entity.Balance{}, err
	}

	b.UserID = id

	return b, nil
}
