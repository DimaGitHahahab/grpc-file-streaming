package client

import (
	"context"
	"google.golang.org/grpc"
	"grpc-file-streaming/internal/service"
	"io"
)

const defaultChunkSize = 256

// Client is a client side of the file service that is used to upload, list, download and remove from Server.
type Client interface {
	Upload(ctx context.Context, fileName string, fileContent []byte) error
	ListFiles(ctx context.Context) ([]string, error)
	Download(ctx context.Context, fileName string) ([]byte, error)
	Delete(ctx context.Context, fileName string) error
}

// client is a gRPC implementation of Client.
type client struct {
	c service.FileServiceClient
}

// New returns a new Client instance with given options.
func New(address string, opts ...grpc.DialOption) (Client, error) {
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		return nil, err
	}
	c := service.NewFileServiceClient(conn)
	return &client{c: c}, nil
}

// Upload sends file name and file content to the serv to be stored.
func (c client) Upload(ctx context.Context, fileName string, fileContent []byte) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		uploadStream, err := c.c.Upload(ctx)
		if err != nil {
			return err
		}
		err = uploadStream.Send(&service.UploadRequest{Content: &service.UploadRequest_FileName{FileName: fileName}})
		if err != nil && err != io.EOF {
			return err
		}
		chunkSize := defaultChunkSize
		for i := 0; i < len(fileContent); i += chunkSize {
			end := i + chunkSize
			if end > len(fileContent) {
				end = len(fileContent)
			}
			err = uploadStream.Send(&service.UploadRequest{Content: &service.UploadRequest_Chunk{Chunk: fileContent[i:end]}})
			if err != nil && err != io.EOF {
				return err
			}
		}
		_, err = uploadStream.CloseAndRecv()
		if err != nil && err != io.EOF {
			return err
		}
	}
	return nil
}

// ListFiles returns all files' meta information from the serv.
func (c client) ListFiles(ctx context.Context) ([]string, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		reply, err := c.c.ListFiles(ctx, &service.ListFilesRequest{})
		if err != nil {
			return nil, err
		}
		return formatMetaData(reply.GetFiles()), nil
	}
}

// Download returns file content from the serv.
func (c client) Download(ctx context.Context, fileName string) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		stream, err := c.c.Download(ctx, &service.DownloadRequest{Name: fileName})
		if err != nil {
			return nil, err
		}
		downloadedContent := make([]byte, 0)
		for {
			downloadResp, err := stream.Recv()
			if err == nil {
				downloadedContent = append(downloadedContent, downloadResp.GetChunk()...)
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}
		}
		return downloadedContent, nil
	}
}

// Delete removes file from the serv.
func (c client) Delete(ctx context.Context, fileName string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		_, err := c.c.Delete(ctx, &service.DeleteRequest{Name: fileName})
		if err != nil {
			return err
		}
	}
	return nil
}
