package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
	"user-balance/domain"
	"user-balance/domain/entity"
	"user-balance/infrastructure/redisdb"
	"user-balance/store"
)

type BalanceService struct {
	store   *store.Store
	redisDB *redisdb.RedisStore
	apiKey  string
}

const converterURL = "https://free.currconv.com/api/v7/"

func NewBalanceService(s *store.Store, rdb *redisdb.RedisStore, ak string) *BalanceService {
	bs := BalanceService{
		store:   s,
		redisDB: rdb,
		apiKey:  ak,
	}

	return &bs
}

func (bs *BalanceService) TransferMoney(t entity.Transaction) (entity.Transaction, error) {
	var isPlus = t.Event == entity.EventCrediting

	t.CreatedAt = time.Now().UTC()

	if t.Event == entity.EventTransfer {
		err := bs.transferringFunds(t)
		if err != nil {
			return entity.Transaction{}, err
		}

		return t, nil
	}

	curB, err := bs.store.SelectUserBalanceByID(t.UserID)
	isNotFound := errors.Is(err, domain.ErrNotFound)

	if isNotFound && !isPlus {
		return entity.Transaction{}, err
	}

	if isNotFound && isPlus {
		err := bs.store.InsertUserTransactions(t)
		if err != nil {
			return entity.Transaction{}, err
		}

		return t, nil
	}

	if err != nil {
		return entity.Transaction{}, err
	}

	if !isPlus && (curB.Amount-t.Amount) < 0 {
		return entity.Transaction{}, domain.ErrEnoughMoney
	}

	if !isPlus {
		t.Amount = -t.Amount
	}

	t.Amount, err = bs.store.UpdateUserTransactions(t)
	if err != nil {
		return entity.Transaction{}, err
	}

	return t, nil
}

func (bs *BalanceService) transferringFunds(t entity.Transaction) error {
	_, err := bs.store.SelectUserBalanceByID(t.TransferID)
	isNotFound := errors.Is(err, domain.ErrNotFound)

	if err != nil && isNotFound {
		b := entity.Balance{
			UserID: t.TransferID,
			Amount: 0,
		}

		err := bs.store.InsertUserBalance(b)
		if err != nil {
			return err
		}
	}

	if err != nil && !isNotFound {
		return err
	}

	err = bs.store.TxUpdateUsersBalances(t)
	if err != nil {
		return err
	}

	return nil
}

func (bs BalanceService) Transactions(uID int64) ([]entity.Transaction, error) {
	ts, err := bs.store.SelectUserTransactions(uID)
	if err != nil && errors.Is(err, domain.ErrNotFound) {
		return make([]entity.Transaction, 0), nil
	}

	if err != nil {
		return nil, domain.ErrUnavailable
	}

	return ts, nil
}

func (bs *BalanceService) UserBalance(id int64, currency string) (entity.Balance, error) {
	b, err := bs.store.SelectUserBalanceByID(id)
	if err != nil {
		return entity.Balance{}, err
	}

	if currency == "" {
		return b, nil
	}

	r, err := http.Get(converterURL + "convert?q=RUB_USD&compact=ultra&apiKey=" + bs.apiKey)
	if err != nil {
		return entity.Balance{}, domain.ErrUnavailable
	}

	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return entity.Balance{}, domain.ErrUnavailable
	}

	var data = make(map[string]interface{})

	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return entity.Balance{}, domain.ErrUnavailable
	}

	v, ok := data["RUB_USD"]
	if !ok {
		return entity.Balance{}, domain.ErrUnavailable
	}

	curs, ok := v.(float64)
	if !ok {
		return entity.Balance{}, domain.ErrUnavailable
	}

	b.Amount = int64(float64(b.Amount) * curs)

	return b, nil
}
