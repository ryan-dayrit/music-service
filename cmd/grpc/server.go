package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"music-service/gen/pb"
	"music-service/internal/config"
	handler "music-service/internal/handler/grpc"
	"music-service/internal/repository/postgres"
	"music-service/pkg/postgres/db"
)

func NewServerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "starts the gRPC server",
		Long:  `starts the gRPC server which hosts MusicService which returns albums`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.Load()
			if err != nil {
				log.Fatalf("failed to load config %v", err)
				return
			}

			address := fmt.Sprintf(":%s", cfg.Grpc.Port)
			listener, err := net.Listen(cfg.Grpc.Network, address)
			if err != nil {
				panic(err)
			}

			s := grpc.NewServer()
			reflection.Register(s)

			db, err := db.NewPostgresDB(cfg.Database)
			if err != nil {
				log.Fatalf("failed to get db: %v", err)
				return
			}
			defer db.Close()

			repository := postgres.NewRepository(db)
			handler := handler.NewHandler(repository)

			pb.RegisterMusicServiceServer(s, handler)
			if err := s.Serve(listener); err != nil {
				log.Fatalf("failed to serve: %v", err)
			}
		},
	}
}
