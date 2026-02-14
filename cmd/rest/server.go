package rest

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"

	"music-service/internal/config"
	"music-service/internal/handler/kafka"
	"music-service/internal/routes"
	v1 "music-service/internal/routes/v1"
	"music-service/pkg/rest"
)

func NewRestServerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "rest_server",
		Short: "starts the REST server",
		Long:  `starts the REST server which hosts MusicService which receives requests to create albums`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.Load()
			if err != nil {
				log.Fatalf("failed to load config %v", err)
				return
			}

			fiberCfg := fiber.Config{
				ReadTimeout:  time.Duration(cfg.Rest.ReadTimeout) * time.Second,
				WriteTimeout: time.Duration(cfg.Rest.WriteTimeout) * time.Second,
			}

			producer, err := kafka.NewProducer(cfg.Kafka)
			if err != nil {
				log.Panicf("Error creating Kafka producer: %v", err)
			}

			app := fiber.New(fiberCfg)

			routes.RegisterSwaggerRoute(app)
			routes.RegisterNotFoundRoute(app)

			v1Router := app.Group("/api/v1")
			v1.RegisterHealthRoute(v1Router)
			v1.RegisterPublicRoutes(v1Router, producer)

			rest.StartServer(app, cfg.Rest)
		},
	}
}
