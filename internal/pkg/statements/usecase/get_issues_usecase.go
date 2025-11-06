package usecase

import (
	"boilerplate-go/internal/pkg/statements/domain/entity"
)

type GetIssuesUsecase interface {
	Execute(uploadID string, statuses []string, page, size int) ([]entity.Transaction, int)
}
