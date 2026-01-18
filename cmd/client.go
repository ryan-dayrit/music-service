package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ryan-dayrit/music-service/config"
	"github.com/ryan-dayrit/music-service/gen/pb"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "shows the albums returned from the MusicService gRPC server",
	Long:  `calls the MusicService gRPC server and shows the albums returned`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			log.Fatalf("failed to load config %v", err)
			return
		}

		target := fmt.Sprintf("%s:%s", cfg.Service.Host, cfg.Service.Port)
		conn, err := grpc.Dial(
			target,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			log.Fatalf("failed to connect: %v", err)
		}
		defer conn.Close()

		client := pb.NewMusicServiceClient(conn)
		getAlbumsResponse, err := client.GetAlbumList(context.Background(), &pb.GetAlbumsRequest{})
		if err != nil {
			log.Fatalf("failed to get album list: %v", err)
		}

		for _, v := range getAlbumsResponse.Albums {
			log.Printf("%v\n", v)
		}
	},
}
