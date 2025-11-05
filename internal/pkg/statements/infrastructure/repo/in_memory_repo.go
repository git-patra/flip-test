package repo

import (
	"boilerplate-go/internal/pkg/statements/domain/entity"
	"sync"
)

type InMemoryRepo struct {
	mu       sync.RWMutex
	byUpload map[string][]entity.Transaction
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{byUpload: map[string][]entity.Transaction{}}
}

func (r *InMemoryRepo) Save(tx entity.Transaction) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byUpload[tx.UploadID] = append(r.byUpload[tx.UploadID], tx)
}

func (r *InMemoryRepo) GetByUpload(uploadID string) []entity.Transaction {
	r.mu.RLock()
	defer r.mu.RUnlock()
	copied := append([]entity.Transaction(nil), r.byUpload[uploadID]...)
	return copied
}
