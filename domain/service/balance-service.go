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
	if b.TypeOp == entity.Plus {
		_, err := bs.store.SelectUserByID(b.UserID)
		if err == domain.ErrNotFound {

			err := bs.store.InsertUserBalance(b)
			if err != nil {
				return entity.Balance{}, err
			}

			return b, nil
		}

		if err != nil {
			return entity.Balance{}, err
		}

		b.Amount, err = bs.store.UpdateUserBalance(b.UserID, b.Amount)
		if err != nil {
			return entity.Balance{}, err
		}

		return b, nil
	}

	if b.TypeOp == entity.Minus {
		curB, err := bs.store.SelectUserByID(b.UserID)
		if errors.Is(err, domain.ErrNotFound) {
			return entity.Balance{}, err
		}

		if err != nil {
			return entity.Balance{}, err
		}

		stateOfB := curB.Amount - b.Amount
		if stateOfB < 0 {
			return entity.Balance{}, domain.ErrEnoughMoney
		}

		b.Amount, err = bs.store.UpdateUserBalance(b.UserID, -b.Amount)
		if err != nil {
			return entity.Balance{}, err
		}

		return b, nil
	}

	return entity.Balance{}, nil
}

func (bs *BalanceService) TransferringFunds(uIdGive, uIdTake, count int) (entity.Balance, error) {
	return entity.Balance{}, nil
}

func (bs *BalanceService) CurrentBalance(id int) (entity.Balance, error) {
	return entity.Balance{}, nil
}
