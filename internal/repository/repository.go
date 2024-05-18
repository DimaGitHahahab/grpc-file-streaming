package repository

import (
	"context"
	"errors"
	"time"

	"grpc-file-streaming/internal/service"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrFileNotFound = status.Error(codes.NotFound, "file not found")

// MongoDBRepository implements server.Repository
type MongoDBRepository struct {
	client   *mongo.Client
	database *mongo.Database
	filesCol *mongo.Collection
}

// NewMongoDBRepository creates a new MongoDBRepository
func NewMongoDBRepository(ctx context.Context, uri, dbName string) (*MongoDBRepository, error) {
	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	database := client.Database(dbName)
	filesCol := database.Collection("files")

	return &MongoDBRepository{
		client:   client,
		database: database,
		filesCol: filesCol,
	}, nil
}

// Get retrieves a file by its name
func (r *MongoDBRepository) Get(ctx context.Context, name string) (*service.File, error) {
	var file service.File
	err := r.filesCol.FindOne(ctx, bson.M{"metadata.name": name}).Decode(&file)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrFileNotFound
		}
		return nil, err
	}
	return &file, nil
}

// Put stores a new file or updates an existing one
func (r *MongoDBRepository) Put(ctx context.Context, name string, content []byte) error {
	filter := bson.M{"metadata.name": name}
	update := bson.M{
		"$set": bson.M{
			"content": content,
			"metadata": service.MetaData{
				Name:      name,
				Size:      int32(len(content)),
				Timestamp: time.Now().Format(time.RFC3339),
			},
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err := r.filesCol.UpdateOne(ctx, filter, update, opts)
	return err
}

// GetAllMetadata retrieves metadata for all files
func (r *MongoDBRepository) GetAllMetadata(ctx context.Context) ([]*service.MetaData, error) {
	cursor, err := r.filesCol.Find(ctx, bson.M{}, options.Find().SetProjection(bson.M{"metadata": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var metaDatas []*service.MetaData
	for cursor.Next(ctx) {
		var result struct {
			Metadata service.MetaData `bson:"metadata"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		metaDatas = append(metaDatas, &result.Metadata)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return metaDatas, nil
}

// Delete removes a file by its name
func (r *MongoDBRepository) Delete(ctx context.Context, name string) error {
	_, err := r.filesCol.DeleteOne(ctx, bson.M{"metadata.name": name})
	return err
}
