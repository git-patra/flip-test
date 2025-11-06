package usecase

import (
	"boilerplate-go/internal/pkg/statements/infrastructure/repo"
)

type getBalanceUsecase struct {
	repo repo.InMemoryRepo
}

func NewGetBalanceUsecase(r repo.InMemoryRepo) GetBalanceUsecase {
	return &getBalanceUsecase{repo: r}
}

func (u *getBalanceUsecase) Execute(uploadID string) int64 {
	txs := u.repo.GetByUpload(uploadID)
	var balance int64
	for _, t := range txs {
		if t.Status == "SUCCESS" {
			if t.Type == "CREDIT" {
				balance += t.Amount
			} else if t.Type == "DEBIT" {
				balance -= t.Amount
			}
		}
	}
	return balance
}
