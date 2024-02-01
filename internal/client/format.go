package client

import (
	"fmt"
	"grpc-file-streaming/internal/service"
)

func formatMetaData(files []*service.MetaData) []string {
	var result []string
	for _, file := range files {
		result = append(result, fmt.Sprintf("name: %v, size: %v, timestamp: %v",
			file.GetName(),
			file.GetSize(),
			file.GetTimestamp(),
		))
	}
	return result
}
