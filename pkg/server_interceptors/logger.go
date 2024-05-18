package server_interceptors

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
)

// StreamLoggerInterceptor logs the beginning and the end of each streaming call:
// duration of execution and error if occurred.
func StreamLoggerInterceptor(srv any,
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	log.Printf("streaming %v started\n", info.FullMethod)

	start := time.Now()

	err := handler(srv, ss)

	l := fmt.Sprintf("streaming of %v completed in %v", info.FullMethod, time.Since(start))
	if err != nil {
		l += fmt.Sprintf(" with error: %v", err)
	}
	log.Println(l)

	return err
}

// UnaryLoggerInterceptor logs the beginning and the end of each unary call:
// duration of execution and error if occurred.
func UnaryLoggerInterceptor(ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	log.Printf("%v started\n", info.FullMethod)

	start := time.Now()

	resp, err := handler(ctx, req)

	l := fmt.Sprintf("%v completed in %v", info.FullMethod, time.Since(start))
	if err != nil {
		l += fmt.Sprintf(" with error: %v", err)
	}

	log.Println(l)

	return resp, err
}
