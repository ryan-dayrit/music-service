package client

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand/v2"
	"net/http"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"music-service/gen/pb"
	"music-service/internal/config"
)

func NewRestClientSingleCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "rest-client-single",
		Short: "sends requests to the MusicService REST server",
		Long:  `calls the MusicService REST server to create an album and shows the response`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.Load()
			if err != nil {
				log.Fatalf("failed to load config %v", err)
			}

			url := "http://" + cfg.Rest.ServerUrl + "/api/v1/album"

			sendAlbumRequest("POST", url)
			sendAlbumRequest("PUT", url)
		},
	}
}

func sendAlbumRequest(method string, url string) {
	album := &pb.Album{
		Id:     rand.Int32(),
		Title:  uuid.NewString(),
		Artist: uuid.NewString(),
		Price:  rand.Float32(),
	}
	jsonData, err := json.Marshal(&album)
	if err != nil {
		log.Fatalf("failed to marshal album %v", err)
	}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("failed to %s album %v", method, err)
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("failed to send %s request %v", method, err)
	}
	defer response.Body.Close()

	log.Printf("%s request response status: %s", method, response.Status)
}
