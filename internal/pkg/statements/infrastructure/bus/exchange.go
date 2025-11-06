package bus

import (
	"context"
	"time"
)

type Exchange interface {
	Publish(exchange string, e Envelope)
	Subscribe(exchange, queue string, filter func(Envelope) bool, buf int) <-chan Envelope
	PublishWithTimeout(ctx context.Context, exchange string, e Envelope, d time.Duration) error
}
