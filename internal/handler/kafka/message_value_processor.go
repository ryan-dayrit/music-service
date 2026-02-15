package kafka

import (
	"log"
	"math/rand/v2"

	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/proto"

	"music-service/gen/pb"
	"music-service/internal/models"
	"music-service/internal/repository/postgres/orm"
)

type MessageValueProcessor struct {
	repository orm.Repository
}

func NewMessageValueProcessor(repository orm.Repository) *MessageValueProcessor {
	return &MessageValueProcessor{repository: repository}
}

func (p *MessageValueProcessor) ProcessMessageValue(messageValue []byte) {
	protoAlbum := &pb.Album{}
	if err := proto.Unmarshal(messageValue, protoAlbum); err != nil {
		log.Fatalf("failed to unmarshal to album: %v", err)
	}

	_, err := p.repository.GetById(int(protoAlbum.Id))
	if err != nil && err.Error() != "pg: no rows in result set" {
		log.Fatalf("failed to read album from postgres: %v", err)
	}

	album := models.Album{
		Id:     int(protoAlbum.Id),
		Title:  protoAlbum.Title,
		Artist: protoAlbum.Artist,
		Price:  decimal.NewFromFloat(rand.Float64()),
	}
	if err != nil && err.Error() == "pg: no rows in result set" {
		err = p.repository.Create(album)
		if err != nil {
			log.Fatalf("failed to create album in postgres: %v", err)
		}
		log.Printf("created album in postgres: %s", album.String())
	} else {
		err = p.repository.Update(album)
		if err != nil {
			log.Fatalf("failed to update album in postgres: %v", err)
		}
		log.Printf("updated album in postgres: %s", album.String())
	}
}
