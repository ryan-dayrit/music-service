package server

import (
	"context"
	"log"

	"music-service/dal/album"
	"music-service/gen/pb"
)

type server struct {
	pb.UnimplementedMusicServiceServer
	album.Repository
}

func NewServer(repository album.Repository) pb.MusicServiceServer {
	return &server{
		Repository: repository,
	}
}

func (s *server) GetAlbumList(context.Context, *pb.GetAlbumsRequest) (*pb.GetAlbumsResponse, error) {
	log.Println("request received")
	albums, err := getAlbumList(s.Repository)
	if err != nil {
		log.Printf("failed to get album list: %v", err)
		return nil, err
	}
	return &pb.GetAlbumsResponse{
		Albums: albums,
	}, nil
}

func getAlbumList(repository album.Repository) ([]*pb.Album, error) {
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
