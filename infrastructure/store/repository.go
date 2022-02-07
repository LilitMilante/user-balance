package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
	"user-balance/domain"
	"user-balance/domain/entity"
	"user-balance/infrastructure/redisdb"
)

type Store struct {
	db  *sql.DB
	rdb *redisdb.RedisStore
}

func NewStore(db *sql.DB, rdb *redisdb.RedisStore) *Store {
	return &Store{
		db:  db,
		rdb: rdb,
	}
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

	_, err = tx.Exec(`INSERT INTO users_transactions (user_id, amount, event, transfer_id, message, created_at) 
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

	rows, err := s.db.Query("SELECT id, user_id, amount, event, transfer_id, message, created_at FROM users_transactions WHERE user_id = $1", uID)
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
	ctx := context.Background()

	userB, err := s.rdb.Get(ctx, strconv.Itoa(int(id))).Bytes()
	if err == nil {
		_ = json.Unmarshal(userB, &b)

		return b, nil
	}

	err = s.db.QueryRow("SELECT user_id, amount FROM users_balances WHERE user_id = $1", id).Scan(&b.UserID, &b.Amount)

	if errors.Is(err, sql.ErrNoRows) {
		return b, domain.ErrNotFound
	}

	if err != nil {
		return b, err
	}

	jb, _ := json.Marshal(b)
	err = s.rdb.Set(ctx, strconv.Itoa(int(id)), jb, 0).Err()
	if err != nil {
		log.Println(err)
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

	_, err = tx.Exec(`INSERT INTO users_transactions (user_id, amount, event, transfer_id, message, created_at) 
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

	err = tx.QueryRow("SELECT amount FROM users_balances WHERE user_id = $1", t.UserID).Scan(&curSum)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.ErrNotFound
		}

		return err
	}

	if (curSum - t.Amount) < 0 {
		return domain.ErrEnoughMoney
	}

	_, err = tx.Exec("UPDATE users_balances SET amount = amount - $1 WHERE user_id = $2", t.Amount, t.UserID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT INTO users_transactions (user_id, amount, event, transfer_id, message, created_at) 
								VALUES ($1, $2, $3, $4, $5, $6)`, t.UserID, t.Amount, entity.EventWriteOffs, t.TransferID, t.Message, t.CreatedAt)
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE users_balances SET amount = amount + $1 WHERE user_id = $2", t.Amount, t.TransferID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT INTO users_transactions (user_id, amount, event, transfer_id, message, created_at) 
								VALUES ($1, $2, $3, $4, $5, $6)`, t.TransferID, t.Amount, entity.EventCrediting, t.UserID, t.Message, t.CreatedAt)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	err = s.rdb.Del(context.Background(), strconv.Itoa(int(t.UserID)), strconv.Itoa(int(t.TransferID))).Err()
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Println("Help Redis!!", err)
	}

	return nil
}
