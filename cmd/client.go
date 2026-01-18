package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "postgres-crud/gen/pb/music"
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
		bookList, err := client.GetAlbumList(context.Background(), &pb.GetAlbumsRequest{})
		if err != nil {
			log.Fatalf("failed to get albums: %v", err)
		}
		log.Printf("albums: %v", bookList)
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
}
