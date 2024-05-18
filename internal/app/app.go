package app

import (
	"context"
	"log"
	"os"

	"grpc-file-streaming/internal/repository"
	"grpc-file-streaming/internal/server"
	"grpc-file-streaming/pkg/config"
	"grpc-file-streaming/pkg/signal"
)

type App struct {
	config  *config.Config
	sigQuit chan os.Signal
	srv     *server.Server
}

func New(cfg *config.Config) *App {
	sigQuit := signal.GetShutdownChannel()

	ctx := context.Background()
	repo, err := repository.NewMongoDBRepository(ctx, cfg.MongoURI, cfg.DbName)
	if err != nil {
		log.Fatal(err)
	}

	s := server.New(repo)

	return &App{
		config:  cfg,
		sigQuit: sigQuit,
		srv:     s,
	}
}

func (a *App) Run() {
	go func() {
		log.Println("Starting server on port:", a.config.Port)

		if err := a.srv.Run(a.config.Port); err != nil {
			log.Fatalln("Failed to start server: ", err)
		}
	}()

	<-a.sigQuit

	log.Println("Gracefully shutting down server")
	a.srv.Stop()
	log.Println("Server shutdown is successful")
}
