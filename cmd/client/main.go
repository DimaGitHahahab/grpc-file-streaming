package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc-file-streaming/internal/client"
	"grpc-file-streaming/pkg/client_interceptors"
	"time"
)

// example of usage of client
func main() {
	c, err := client.New("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(client_interceptors.UnaryLoggerInterceptor),
		grpc.WithStreamInterceptor(client_interceptors.StreamLoggerInterceptor),
	)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = c.Upload(ctx, "test.txt", []byte("test"))
	if err != nil {
		panic(err)
	}

	files, err := c.ListFiles(ctx)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		println(file)
	}

	fileContent, err := c.Download(ctx, "test.txt")
	if err != nil {
		panic(err)
	}
	println(string(fileContent))

	err = c.Delete(ctx, "test.txt")
	if err != nil {
		panic(err)
	}

}
