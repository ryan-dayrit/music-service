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

func (s *server) GetAlbumList(ctx context.Context, req *pb.GetAlbumsRequest) (*pb.GetAlbumsResponse, error) {
	log.Println("request received")
	
	// Check if context is already canceled
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	
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
