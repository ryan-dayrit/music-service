package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "music-service/gen/pb/music"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "shows the albums returned from the MusicService gRPC server",
	Long:  `calls the MusicService gRPC server and shows the albums returned`,
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
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
