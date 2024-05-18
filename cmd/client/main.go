package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"grpc-file-streaming/internal/client"
)

// example of client methods usage
func main() {
	c, err := client.New("localhost:50051")
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// uploading file name and its payload
	err = c.Upload(ctx, "example.txt", []byte("some example payload"))
	if err != nil {
		log.Fatal(err)
	}

	// list of all the files in the server
	files, err := c.ListFiles(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		fmt.Println(file)
	}

	fileContent, err := c.Download(ctx, "example.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(fileContent))

	// removing file from server
	err = c.Delete(ctx, "example.txt")
	if err != nil {
		log.Fatal(err)
	}
}
