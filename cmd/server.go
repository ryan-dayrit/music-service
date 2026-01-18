package cmd

import (
	"log"
	"net"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/ryan-dayrit/music-service/config"
	"github.com/ryan-dayrit/music-service/dal/album"
	"github.com/ryan-dayrit/music-service/db"
	"github.com/ryan-dayrit/music-service/gen/pb"
	"github.com/ryan-dayrit/music-service/server"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "starts the gRPC server",
	Long:  `starts the gRPC server which hosts MusicService which returns albums`,
	Run: func(cmd *cobra.Command, args []string) {
		listener, err := net.Listen("tcp", ":8080")
		if err != nil {
			panic(err)
		}

		s := grpc.NewServer()
		reflection.Register(s)

		cfg, err := config.Load()
		if err != nil {
			log.Fatalf("failed to load config %v", err)
			return
		}

		db, err := db.GetDB(*cfg)
		if err != nil {
			log.Fatalf("failed to get db: %v", err)
			return
		}
		defer db.Close()

		repository := album.NewRepository(db)
		server := server.NewServer(repository)

		pb.RegisterMusicServiceServer(s, server)
		if err := s.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	},
}
