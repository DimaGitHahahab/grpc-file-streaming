package client_interceptors

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"time"
)

// StreamLoggerInterceptor logs the beginning and the end of each streaming call:
// duration of execution and error if occurred.
func StreamLoggerInterceptor(ctx context.Context,
	desc *grpc.StreamDesc,
	cc *grpc.ClientConn,
	method string,
	streamer grpc.Streamer,
	opts ...grpc.CallOption,
) (grpc.ClientStream, error) {
	log.Printf("streaming %v started\n", method)
	start := time.Now()

	clientStream, err := streamer(ctx, desc, cc, method, opts...)

	l := fmt.Sprintf("streaming %v completed in %v", method, time.Since(start))
	if err != nil {
		l += fmt.Sprintf(" with error: %v", err)
	}

	log.Println(l)

	return clientStream, err
}

// UnaryLoggerInterceptor logs the beginning and the end of each unary call:
func UnaryLoggerInterceptor(ctx context.Context,
	method string,
	req,
	reply any, cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	log.Printf("%v started\n", method)

	start := time.Now()

	err := invoker(ctx, method, req, reply, cc, opts...)

	l := fmt.Sprintf("%v completed in %v", method, time.Since(start))
	if err != nil {
		l += fmt.Sprintf(" with error: %v", err)
	}

	log.Println(l)

	return err
}
