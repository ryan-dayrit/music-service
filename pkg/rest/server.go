package rest

import (
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
)

func StartServerWithGracefulShutdown(app *fiber.App, cfg Config) {
	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := app.Shutdown(); err != nil {
			log.Printf("server is not shutting down: %v", err)
		}

		close(idleConnsClosed)
	}()

	StartServer(app, cfg)

	<-idleConnsClosed
}

func StartServer(app *fiber.App, cfg Config) {
	if err := app.Listen(cfg.ServerUrl); err != nil {
		log.Printf("server is not running: %v", err)
	}
}
