package message

import (
	"errors"
	"log"

	"github.com/go-pg/pg/v10"
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

func (p *MessageValueProcessor) Process(messageValue []byte) {
	protoAlbum := &pb.Album{}
	if err := proto.Unmarshal(messageValue, protoAlbum); err != nil {
		log.Printf("skipping message: failed to unmarshal album: %v", err)
		return
	}

	_, err := p.repository.GetById(int(protoAlbum.Id))
	if err != nil && !errors.Is(err, pg.ErrNoRows) {
		log.Printf("skipping message: failed to read album from postgres: %v", err)
		return
	}

	album := models.Album{
		Id:     int(protoAlbum.Id),
		Title:  protoAlbum.Title,
		Artist: protoAlbum.Artist,
		Price:  decimal.NewFromFloat(float64(protoAlbum.Price)),
	}
	if errors.Is(err, pg.ErrNoRows) {
		if err := p.repository.Create(album); err != nil {
			log.Printf("failed to create album in postgres: %v", err)
			return
		}
		log.Printf("created album in postgres: %s", album.String())
	} else {
		if err := p.repository.Update(album); err != nil {
			log.Printf("failed to update album in postgres: %v", err)
			return
		}
		log.Printf("updated album in postgres: %s", album.String())
	}
}
