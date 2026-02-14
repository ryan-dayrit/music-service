package rest

import (
	"log"

	"github.com/spf13/cobra"

	"music-service/internal/config"
)

func NewRestClientCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "rest_client",
		Short: "sends requests to the MusicService REST server",
		Long:  `calls the MusicService REST server to create an album and shows the response`,
		Run: func(cmd *cobra.Command, args []string) {
			_, err := config.Load()
			if err != nil {
				log.Fatalf("failed to load config %v", err)
				return
			}

			// TODO: implement the REST client that sends requests to the REST server and shows the response
		},
	}
}
