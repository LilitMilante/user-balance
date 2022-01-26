package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"user-balance/domain"
	"user-balance/domain/entity"
	"user-balance/store"
)

type BalanceService struct {
	store  *store.Store
	apiKey string
}

const converterURL = "https://free.currconv.com/api/v7/"

func NewBalanceService(s *store.Store, ak string) *BalanceService {
	bs := BalanceService{
		store:  s,
		apiKey: ak,
	}

	return &bs
}

func (bs *BalanceService) BalanceOperations(b entity.Balance) (entity.Balance, error) {
	var isPlus = b.TypeOp == entity.Plus

	curB, err := bs.store.SelectUserBalanceByID(b.UserID)
	isNotFound := errors.Is(err, domain.ErrNotFound)

	if isNotFound && !isPlus {
		return entity.Balance{}, err
	}

	if isNotFound && isPlus {
		err := bs.store.InsertUserTransactions(b)
		if err != nil {
			return entity.Balance{}, err
		}

		return b, nil
	}

	if err != nil {
		return entity.Balance{}, err
	}

	if !isPlus && (curB.Amount-b.Amount) < 0 {
		return entity.Balance{}, domain.ErrEnoughMoney
	}

	if !isPlus {
		b.Amount = -b.Amount
	}

	b.Amount, err = bs.store.UpdateUserTransactions(b)
	if err != nil {
		return entity.Balance{}, err
	}

	return b, nil

}

func (bs *BalanceService) TransferringFunds(t entity.Transfer) error {
	_, err := bs.store.SelectUserBalanceByID(t.IdTake)
	isNotFound := errors.Is(err, domain.ErrNotFound)

	if err != nil && isNotFound {
		b := entity.Balance{
			UserID: t.IdTake,
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

func (bs BalanceService) Transactions(uID int64) ([]entity.Balance, error) {
	ts, err := bs.store.SelectUserTransactions(uID)
	if err != nil && errors.Is(err, domain.ErrNotFound) {
		return make([]entity.Balance, 0), nil
	}

	if err != nil {
		return nil, domain.ErrUnavailable
	}

	return ts, nil
}
