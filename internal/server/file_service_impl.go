package server

import (
	"context"
	"io"

	"grpc-file-streaming/internal/service"
)

const defaultChunkSize = 256

type Repository interface {
	Get(context.Context, string) (*service.File, error)
	Put(context.Context, string, []byte) error
	GetAllMetadata(context.Context) ([]*service.MetaData, error)
	Delete(context.Context, string) error
}

// fileServiceServer is a gRPC implementation of FileServiceServer.
type fileServiceServer struct {
	service.UnimplementedFileServiceServer
	repo Repository
}

// Download accepts service.DownloadRequest, gets file name from it and streams the file back to the client.
func (s *fileServiceServer) Download(request *service.DownloadRequest, stream service.FileService_DownloadServer) error {
	fileName := request.GetName()
	file, err := s.repo.Get(stream.Context(), fileName)
	if err != nil {
		return err
	}

	fileContent := file.GetContent()
	chunkSize := defaultChunkSize
	for i := 0; i < len(fileContent); i += chunkSize {
		end := i + chunkSize
		if end > len(fileContent) {
			end = len(fileContent)
		}
		downloadResponse := &service.DownloadResponse{Chunk: fileContent[i:end]}
		if err = stream.Send(downloadResponse); err != nil {
			return err
		}
	}
	return nil
}

// ListFiles sends all files' descriptions (service.MetaData) without their content to the client.
func (s *fileServiceServer) ListFiles(ctx context.Context, _ *service.ListFilesRequest) (*service.ListFilesResponse, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			metaData, err := s.repo.GetAllMetadata(ctx)
			if err != nil {
				return nil, err
			}
			return &service.ListFilesResponse{Files: metaData}, nil
		}
	}
}

// Upload accepts service.UploadRequest, gets file name or content from every sent portion of data
// and then stores the file in the Repository.
func (s *fileServiceServer) Upload(stream service.FileService_UploadServer) error {
	var fileName string
	var fileContent []byte
	for {
		uploadRequest, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		switch content := uploadRequest.Content.(type) {
		case *service.UploadRequest_Chunk:
			fileContent = append(fileContent, content.Chunk...)
		case *service.UploadRequest_FileName:
			fileName = content.FileName

		}
	}
	return s.repo.Put(stream.Context(), fileName, fileContent)
}

// Delete accepts service.DeleteRequest, gets file name from it and deletes the file from the Repository.
func (s *fileServiceServer) Delete(ctx context.Context, in *service.DeleteRequest) (*service.DeleteResponse, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			err := s.repo.Delete(ctx, in.GetName())
			if err != nil {
				return nil, err
			}
			return &service.DeleteResponse{}, nil
		}
	}
}
