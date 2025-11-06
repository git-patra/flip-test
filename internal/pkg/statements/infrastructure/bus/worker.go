package bus

import "context"

type Worker interface {
	Start(ctx context.Context)
}
