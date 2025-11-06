package usecase

type GetBalanceUsecase interface {
	Execute(uploadID string) int64
}
