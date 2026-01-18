package cmd

import (
	"context"
	"log"
	"net"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "postgres-crud/gen/pb/music"
)

type server struct {
	pb.UnimplementedMusicServiceServer
}

func (s *server) GetAlbums(context.Context, *pb.GetAlbumsRequest) (*pb.GetAlbumsResponse, error) {
	return &pb.GetAlbumsResponse{
		Albums: getAlbums(),
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
			log.Fatalf("Failed to serve: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func getAlbums() []*pb.Album {
	albums := []*pb.Album{
		{
			Id:     1,
			Title:  "Blue Train'",
			Artist: "John Coltrane",
			Price:  56.99,
		},
	}
	return albums
}
