package repo

import (
	"boilerplate-go/internal/pkg/statements/domain/entity"
)

type InMemoryRepo interface {
	Save(tx entity.Transaction)
	GetByUpload(uploadID string) []entity.Transaction
}
