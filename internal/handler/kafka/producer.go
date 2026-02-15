package kafka

import (
	"context"

	"music-service/gen/pb"
)

type ProducerHandler interface {
	Produce(ctx context.Context, album *pb.Album)
}
