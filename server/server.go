package server

import (
	"context"
	"log"

	"github.com/ryan-dayrit/music-service/dal/album"
	"github.com/ryan-dayrit/music-service/gen/pb"
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
	return &pb.GetAlbumsResponse{
		Albums: getAlbumList(s.Repository),
	}, nil
}

func getAlbumList(repository album.Repository) []*pb.Album {
	albums, err := repository.Read()
	if err != nil {
		log.Fatalf("failed to read albums: %v", err)
		return nil
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
	return albumList
}
