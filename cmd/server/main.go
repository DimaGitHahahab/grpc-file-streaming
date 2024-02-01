package main

import (
	"google.golang.org/grpc"
	"grpc-file-streaming/internal/repository"
	"grpc-file-streaming/internal/server"
	"grpc-file-streaming/pkg/server_interceptors"
	"log"
)

// example of usage of server
func main() {

	unaryInterceptorChain := grpc.ChainUnaryInterceptor(
		server_interceptors.UnaryValidate,
		server_interceptors.UnaryLoggerInterceptor,
	)

	streamInterceptorChain := grpc.ChainStreamInterceptor(
		server_interceptors.StreamValidate,
		server_interceptors.StreamLoggerInterceptor,
	)

	s := server.New(repository.New(), unaryInterceptorChain, streamInterceptorChain)
	if err := s.Listen(":50051"); err != nil {
		log.Fatal("failed to start server: ", err)
	}
}
