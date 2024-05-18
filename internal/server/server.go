package server

import (
	"net"

	"grpc-file-streaming/internal/service"
	"grpc-file-streaming/pkg/server_interceptors"

	"google.golang.org/grpc"
)

// New returns a new Server with interceptors.
// To start the server, call Listen method.
func New(repo Repository) *Server {
	unaryInterceptorChain := grpc.ChainUnaryInterceptor(
		server_interceptors.UnaryValidate,
		server_interceptors.UnaryLoggerInterceptor,
	)

	streamInterceptorChain := grpc.ChainStreamInterceptor(
		server_interceptors.StreamValidate,
		server_interceptors.StreamLoggerInterceptor,
	)

	s := grpc.NewServer(unaryInterceptorChain, streamInterceptorChain)

	service.RegisterFileServiceServer(s, &fileServiceServer{repo: repo})

	return &Server{srv: s}
}

type Server struct {
	srv *grpc.Server
}

func (s *Server) Run(addr string) error {
	lis, err := net.Listen("tcp", net.JoinHostPort("", addr))
	if err != nil {
		return err
	}
	return s.srv.Serve(lis)
}

func (s *Server) Stop() {
	s.srv.GracefulStop()
}
