package client

import (
	"context"
	"testing"
	"time"

	"grpc-file-streaming/internal/server"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestClient_ContextCancelled(t *testing.T) {
	serv := server.New(nil)
	go func() {
		_ = serv.Listen(":50051")
		time.Sleep(10 * time.Second)
	}()

	c, _ := New(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	cancel()

	err := c.Upload(ctx, "test.txt", []byte("test"))
	assert.Error(t, err)

	_, err = c.Download(ctx, "test.txt")
	assert.Error(t, err)

	err = c.Delete(ctx, "test.txt")
	assert.Error(t, err)

	_, err = c.ListFiles(ctx)
	assert.Error(t, err)
}

func TestClient_ContextTimeout(t *testing.T) {
	serv := server.New(nil)
	go func() {
		_ = serv.Listen(":50051")
		time.Sleep(10 * time.Second)
	}()

	c, _ := New(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	err := c.Upload(ctx, "test.txt", []byte("test"))
	assert.Error(t, err)

	_, err = c.Download(ctx, "test.txt")
	assert.Error(t, err)

	err = c.Delete(ctx, "test.txt")
	assert.Error(t, err)

	_, err = c.ListFiles(ctx)
	assert.Error(t, err)
}
