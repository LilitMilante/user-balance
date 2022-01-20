package service

import (
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

func (bs *BalanceService) CreditingFunds(b entity.Balance) (entity.Balance, error) {
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

	b, err = bs.store.UpdateUserBalance(b.UserID, b.Amount)
	if err != nil {
		return entity.Balance{}, err
	}

	return b, nil
}

func (bs *BalanceService) WriteOffsFunds(uId, count int) (entity.Balance, error) {
	return entity.Balance{}, nil
}

func (bs *BalanceService) TransferringFunds(uIdGive, uIdTake, count int) (entity.Balance, error) {
	return entity.Balance{}, nil
}

func (bs *BalanceService) CurrentBalance(id int) (entity.Balance, error) {
	return entity.Balance{}, nil
}
