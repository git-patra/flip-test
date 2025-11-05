package bus

import (
	"boilerplate-go/internal/pkg/statements/domain/entity"
	"context"
	"log"
	"sync"
	"time"
)

type Handler func(ctx context.Context, evt entity.FailedTransactionOccurred) error

type WorkerPool struct {
	workers, maxRetries int
	in                  <-chan entity.FailedTransactionOccurred
	handler             Handler
	processed           sync.Map
}

func NewWorkerPool(workers, maxRetries int, in <-chan entity.FailedTransactionOccurred, handler Handler) *WorkerPool {
	return &WorkerPool{workers: workers, maxRetries: maxRetries, in: in, handler: handler}
}

func (w *WorkerPool) Start(ctx context.Context) {
	wg := sync.WaitGroup{}
	for i := 0; i < w.workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case evt := <-w.in:
					key := evt.UploadID + "/" + evt.Transaction.ID
					if _, seen := w.processed.Load(key); seen {
						continue
					}
					backoff := 200 * time.Millisecond
					var err error
					for attempt := 0; attempt <= w.maxRetries; attempt++ {
						err = w.handler(ctx, evt)
						if err == nil {
							w.processed.Store(key, struct{}{})
							break
						}
						select {
						case <-ctx.Done():
							return
						case <-time.After(backoff):
							backoff *= 2
						}
					}
					if err != nil {
						log.Printf("worker[%d]: failed after retries upload=%s tx=%s err=%v", id, evt.UploadID, evt.Transaction.ID, err)
					}
				}
			}
		}(i)
	}
	go func() { wg.Wait() }()
}
