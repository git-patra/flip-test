package bus

import (
	"context"
	"log"
	"sync"
	"time"
)

// ===== Envelope & exchange
type Envelope struct {
	RoutingKey string      // e.g. "transactions.failed"
	Key        string      // idempotency key, e.g. "uploadID/txID"
	Payload    interface{} // concrete domain payload
}

type Subscriber struct {
	QueueName string
	Ch        chan Envelope
	Match     func(Envelope) bool // binding filter (routing key/topic)
}

type exchange struct {
	mu   sync.RWMutex
	subs map[string][]*Subscriber // exchangeName -> subscribers
}

func NewExchange() Exchange {
	return &exchange{subs: map[string][]*Subscriber{}}
}

func (x *exchange) Publish(exchange string, e Envelope) {
	x.mu.RLock()
	defer x.mu.RUnlock()
	list := x.subs[exchange]
	if len(list) == 0 {
		log.Printf("[WARN] no subscribers for exchange=%s rk=%s", exchange, e.RoutingKey)
		return
	}
	for _, s := range list {
		if s.Match == nil || s.Match(e) {
			select {
			case s.Ch <- e:
			default:
				log.Printf("[WARN] queue=%s full, drop rk=%s key=%s", s.QueueName, e.RoutingKey, e.Key)
			}
		}
	}
}

func (x *exchange) Subscribe(exchange, queue string, filter func(Envelope) bool, buf int) <-chan Envelope {
	x.mu.Lock()
	defer x.mu.Unlock()
	s := &Subscriber{
		QueueName: queue,
		Ch:        make(chan Envelope, buf),
		Match:     filter,
	}
	x.subs[exchange] = append(x.subs[exchange], s)
	return s.Ch
}

// Optional blocking/timeout publish
func (x *exchange) PublishWithTimeout(ctx context.Context, exchange string, e Envelope, d time.Duration) error {
	done := make(chan struct{}, 1)
	go func() { x.Publish(exchange, e); done <- struct{}{} }()
	select {
	case <-done:
		return nil
	case <-time.After(d):
		return context.DeadlineExceeded
	case <-ctx.Done():
		return ctx.Err()
	}
}
