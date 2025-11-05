package bus

import (
	"context"
	"log"
	"sync"
	"time"
)

type HandlerFn func(ctx context.Context, env Envelope) error

type Consumer struct {
	workers    int
	maxRetries int
	in         <-chan Envelope
	handler    HandlerFn
	processed  sync.Map // idempotency by Envelope.Key
}

func NewConsumer(workers, maxRetries int, in <-chan Envelope, h HandlerFn) *Consumer {
	if workers <= 0 {
		workers = 4
	}
	if maxRetries < 0 {
		maxRetries = 0
	}
	return &Consumer{workers: workers, maxRetries: maxRetries, in: in, handler: h}
}

func (c *Consumer) Start(ctx context.Context) {
	wg := sync.WaitGroup{}
	for i := 0; i < c.workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case env := <-c.in:
					if _, seen := c.processed.Load(env.Key); seen {
						continue
					}
					backoff := 200 * time.Millisecond
					var err error
					for attempt := 0; attempt <= c.maxRetries; attempt++ {
						err = c.handler(ctx, env)
						if err == nil {
							c.processed.Store(env.Key, struct{}{})
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
						log.Printf(`[ERROR] consumer[%d] rk=%s key=%s err=%v`, id, env.RoutingKey, env.Key, err)
					}
				}
			}
		}(i)
	}
	go func() { wg.Wait() }()
}
