package bus

import (
	"boilerplate-go/internal/pkg/statements/domain/entity"
	"context"
	"time"
)

type InMemoryBus struct {
	ch chan entity.FailedTransactionOccurred
}

func NewInMemoryBus(buf int) *InMemoryBus {
	if buf <= 0 {
		buf = 1024
	}
	return &InMemoryBus{ch: make(chan entity.FailedTransactionOccurred, buf)}
}

func (b *InMemoryBus) PublishNonBlocking(evt entity.FailedTransactionOccurred) {
	select {
	case b.ch <- evt:
	default:
		// drop event jika penuh
	}
}

func (b *InMemoryBus) PublishWithTimeout(ctx context.Context, evt entity.FailedTransactionOccurred, d time.Duration) error {
	select {
	case b.ch <- evt:
		return nil
	case <-time.After(d):
		return context.DeadlineExceeded
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (b *InMemoryBus) Subscribe() <-chan entity.FailedTransactionOccurred {
	return b.ch
}
