package main

import (
	"grpc-file-streaming/internal/app"
	"grpc-file-streaming/pkg/config"
)

func main() {
	cfg := config.LoadConfig()

	a := app.New(cfg)

	a.Run()
}
