package repo

import (
	"boilerplate-go/internal/pkg/statements/domain/entity"
	"sync"
)

type inMemoryRepo struct {
	mu       sync.RWMutex
	byUpload map[string][]entity.Transaction
}

func NewInMemoryRepo() InMemoryRepo {
	return &inMemoryRepo{byUpload: map[string][]entity.Transaction{}}
}

func (r *inMemoryRepo) Save(tx entity.Transaction) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byUpload[tx.UploadID] = append(r.byUpload[tx.UploadID], tx)
}

func (r *inMemoryRepo) GetByUpload(uploadID string) []entity.Transaction {
	r.mu.RLock()
	defer r.mu.RUnlock()
	copied := append([]entity.Transaction(nil), r.byUpload[uploadID]...)
	return copied
}
