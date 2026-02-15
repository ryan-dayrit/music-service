package client

import (
	"bytes"
	"log"
	"net/http"

	"encoding/json"
	"math/rand/v2"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"music-service/gen/pb"
	"music-service/internal/config"
)

func NewRestClientMultiCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "rest-client-multi",
		Short: "sends requests to the MusicService REST server",
		Long:  `calls the MusicService REST server to create multiple albums and shows the response`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.Load()
			if err != nil {
				log.Fatalf("failed to load config %v", err)
			}

			url := "http://" + cfg.Rest.ServerUrl + "/api/v1/albums"

			sendAlbumsRequest("POST", url)
			sendAlbumsRequest("PUT", url)
		},
	}
}

func sendAlbumsRequest(method string, url string) {
	albums := []*pb.Album{}

	for i := 0; i < 10; i++ {
		album := &pb.Album{
			Id:     rand.Int32(),
			Title:  uuid.NewString(),
			Artist: uuid.NewString(),
			Price:  rand.Float32(),
		}
		albums = append(albums, album)
	}

	jsonData, err := json.Marshal(&albums)
	if err != nil {
		log.Fatalf("failed to marshal albums %v", err)
	}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("failed to %s albums %v", method, err)
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
