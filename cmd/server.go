package cmd

import (
	"context"
	"log"
	"net"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/ryan-dayrit/music-service/config"
	"github.com/ryan-dayrit/music-service/dal/album"
	"github.com/ryan-dayrit/music-service/db"
	pb "github.com/ryan-dayrit/music-service/gen/pb/music"
)

type server struct {
	pb.UnimplementedMusicServiceServer
}

func (s *server) GetAlbumList(context.Context, *pb.GetAlbumsRequest) (*pb.GetAlbumsResponse, error) {
	log.Println("request received")
	return &pb.GetAlbumsResponse{
		Albums: getAlbumList(),
	}, nil
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "starts the MusicService gRPC server",
	Long:  `starts the MusicService gRPC which returns albums`,
	Run: func(cmd *cobra.Command, args []string) {
		listener, err := net.Listen("tcp", ":8080")
		if err != nil {
			panic(err)
		}

		s := grpc.NewServer()
		reflection.Register(s)
		pb.RegisterMusicServiceServer(s, &server{})
		if err := s.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	},
}

func getAlbumList() []*pb.Album {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config %v", err)
		return nil
	}

	db, err := db.GetDB(*cfg)
	if err != nil {
		log.Fatalf("failed to get db: %v", err)
		return nil
	}
	defer db.Close()

	repository := album.NewRepository(db)

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
