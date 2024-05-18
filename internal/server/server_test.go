package server

import (
	"context"
	"testing"
	"time"

	"grpc-file-streaming/internal/client"
	"grpc-file-streaming/internal/repository"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestUpload(t *testing.T) {
	st := repository.New()
	serv := New(st)

	go func() {
		_ = serv.Listen(":50051")
		time.Sleep(10 * time.Second)
	}()

	c, err := client.New(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = c.Upload(ctx, "test.txt", []byte("test"))
	assert.NoError(t, err)

	file, err := st.Get("test.txt")
	assert.NoError(t, err)
	assert.Equal(t, "test", string(file.Content))
}

func TestDownload(t *testing.T) {
	st := repository.New()
	_ = st.Put("test.txt", []byte("test"))
	serv := New(st)

	go func() {
		_ = serv.Listen(":50051")
		time.Sleep(10 * time.Second)
	}()

	c, err := client.New(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	content, err := c.Download(ctx, "test.txt")
	assert.NoError(t, err)
	assert.Equal(t, "test", string(content))
}

func TestDelete(t *testing.T) {
	st := repository.New()
	_ = st.Put("test.txt", []byte("test"))
	serv := New(st)

	go func() {
		_ = serv.Listen(":50051")
		time.Sleep(10 * time.Second)
	}()

	c, err := client.New(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = c.Delete(ctx, "test.txt")
	assert.NoError(t, err)
}

func TestListFiles(t *testing.T) {
	st := repository.New()
	_ = st.Put("test.txt", []byte("test"))
	serv := New(st)

	go func() {
		_ = serv.Listen(":50051")
		time.Sleep(10 * time.Second)
	}()

	c, err := client.New(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	metadata, err := c.ListFiles(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, "test.txt", metadata)
}
