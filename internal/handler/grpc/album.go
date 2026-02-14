package grpc

import (
	"context"
	"log"

	"music-service/gen/pb"
	"music-service/internal/repository/postgres"
)

type albumHandler struct {
	pb.UnimplementedMusicServiceServer
	postgres.Repository
}

func NewAlbumHandler(repository postgres.Repository) pb.MusicServiceServer {
	return &albumHandler{
		Repository: repository,
	}
}

func (h *albumHandler) GetAlbumList(ctx context.Context, req *pb.GetAlbumsRequest) (*pb.GetAlbumsResponse, error) {
	log.Println("request received")

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	albums, err := getAlbumList(h.Repository)
	if err != nil {
		return nil, err
	}

	return &pb.GetAlbumsResponse{
		Albums: albums,
	}, nil
}

func getAlbumList(repository postgres.Repository) ([]*pb.Album, error) {
	albums, err := repository.Read()
	if err != nil {
		return nil, err
	}

	albumList := make([]*pb.Album, len(albums))
	for i, v := range albums {
		priceF64, _ := v.Price.Float64()
		albumList[i] = &pb.Album{
			Id:     int32(v.Id),
			Title:  v.Title,
			Artist: v.Artist,
			Price:  float32(priceF64),
		}
	}
	return albumList, nil
}
