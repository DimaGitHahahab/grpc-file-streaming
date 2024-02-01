package repository

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc-file-streaming/internal/service"
	"sync"
	"time"
)

var (
	ErrFileNotFound = status.Error(codes.NotFound, "file not found")
	ErrFileExists   = status.Error(codes.AlreadyExists, "file already exists")
)

// FileStorage is an in-memory concurrent-safe implementation of server.Repository.
type FileStorage struct {
	files map[string]*service.File
	mu    sync.RWMutex
}

// New returns a new FileStorage instance.
func New() *FileStorage {
	return &FileStorage{files: make(map[string]*service.File), mu: sync.RWMutex{}}
}

// Get returns a file by name.
func (s *FileStorage) Get(name string) (*service.File, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if file, ok := s.files[name]; ok {
		return file, nil
	}
	return nil, ErrFileNotFound
}

// Put adds a new file to the repository.
func (s *FileStorage) Put(name string, content []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.files[name]; ok {
		return ErrFileExists
	}
	metaData := &service.MetaData{Name: name, Size: int32(len(content)), Timestamp: time.Now().String()}
	file := &service.File{Content: content, Metadata: metaData}
	s.files[file.GetMetadata().GetName()] = file
	return nil
}

// GetAllMetadata returns metadata of all files in the repository.
func (s *FileStorage) GetAllMetadata() ([]*service.MetaData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]*service.MetaData, 0, len(s.files))
	for _, file := range s.files {
		res = append(res, file.GetMetadata())
	}
	return res, nil
}

// Delete removes a file from the repository by name.
func (s *FileStorage) Delete(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.files[name]; !ok {
		return ErrFileNotFound
	}
	delete(s.files, name)
	return nil
}
