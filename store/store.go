package store

import (
	"io"

	"github.com/alec-rabold/copypasta/runcommands"

	minio "github.com/minio/minio-go"
	log "github.com/sirupsen/logrus"
)

// Store is the interface for reading from / writing to a storage destination
type Store interface {
	Write(content io.Reader) error
	Read() (io.Reader, error)
}

// NewStore creates new s3 destination client
func NewStore(target *runcommands.S3Target) (Store, error) {
	client, err := minioClient(target)
	if err != nil {
		log.Errorf("Error initializing client: %s", err.Error())
		return nil, err
	}
	return NewS3Store(client, target), nil
}

func minioClient(t *runcommands.S3Target) (*minio.Client, error) {
	useSSL := true
	return minio.New(t.Endpoint, t.AccessKey, t.SecretAccessKey, useSSL)
}
