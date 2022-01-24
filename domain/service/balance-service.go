package service

import (
	"errors"
	"user-balance/domain"
	"user-balance/domain/entity"
	"user-balance/store"
)

type BalanceService struct {
	store *store.Store
}

func NewBalanceService(s *store.Store) *BalanceService {
	bs := BalanceService{
		store: s,
	}

	return &bs
}

func (bs *BalanceService) BalanceOperations(b entity.Balance) (entity.Balance, error) {
	var isPlus = b.TypeOp == entity.Plus

	curB, err := bs.store.SelectUserByID(b.UserID)
	isNotFound := errors.Is(err, domain.ErrNotFound)

	if isNotFound && !isPlus {
		return entity.Balance{}, err
	}

	if isNotFound && isPlus {
		err := bs.store.InsertUserBalance(b)
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

	b.Amount, err = bs.store.UpdateUserBalance(b.UserID, b.Amount)
	if err != nil {
		return entity.Balance{}, err
	}

	return b, nil

}

func (bs *BalanceService) TransferringFunds(uIdGive, uIdTake, count int) (entity.Balance, error) {
	return entity.Balance{}, nil
}

func (bs *BalanceService) CurrentBalance(id int) (entity.Balance, error) {
	return entity.Balance{}, nil
}
