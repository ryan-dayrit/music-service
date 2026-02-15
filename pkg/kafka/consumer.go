package kafka

import (
	"context"
)

type ConsumerHandler interface {
	Consume(ctx context.Context) error
}
