package store

import (
	"bytes"
	"github.com/alec-rabold/copypasta/runcommands"
	minio "github.com/minio/minio-go"
	"io"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
)

const (
	defaultObjectName        = "default-object-name"
	defaultContentType       = "text/html"
	defaultStreamSizeUnknown = -1
)

// MinioClient is the interface to read and write to s3
type MinioClient interface {
	MakeBucket(string, string) error
	BucketExists(string) (bool, error)
	PutObject(string, string, io.Reader, int64, minio.PutObjectOptions) (int64, error)
	FGetObject(string, string, string, minio.GetObjectOptions) error
}

// S3Store is the struct for client and destination
type S3Store struct {
	minioClient MinioClient
	target      *runcommands.S3Target
}

// NewS3Store creates new client and target
func NewS3Store(client MinioClient, target *runcommands.S3Target) *S3Store {
	return &S3Store{
		minioClient: client,
		target:      target,
	}
}

// TODO: Encryption option/flag
// Write is the function responsible for writing to s3
func (s *S3Store) Write(content io.Reader) error {
	bucketExists, err := s.minioClient.BucketExists(s.target.BucketName)
	if err != nil {
		log.Errorf("Error checking for bucket: %s", err.Error())
		return err
	}

	if !bucketExists {
		err := s.minioClient.MakeBucket(s.target.BucketName, s.target.Location)
		if err != nil {
			log.Errorf("Error creating bucket: %s", err.Error())
			return err
		}
	}

	_, err = s.minioClient.PutObject(s.target.BucketName, defaultObjectName, content, defaultStreamSizeUnknown, minio.PutObjectOptions{ContentType: defaultContentType})
	if err != nil {
		log.Errorf("Error writing to bucket: %s", err.Error())
		return err
	}

	return nil
}

// Read is the function responsible to reading from s3
func (s *S3Store) Read() (io.Reader, error) {
	tempFile, err := ioutil.TempFile(os.TempDir(), "tempS3ObjectFile")
	if err != nil {
		log.Errorf("Error reading from bucket: %s", err.Error())
		return nil, err
	}
	defer func() {
		tempFile.Close()
		if err = os.Remove(tempFile.Name()); err != nil {
			log.Fatal("Error removing temp s3 file", err.Error())
		}
	}()
	err = s.minioClient.FGetObject(s.target.BucketName, defaultObjectName, tempFile.Name(), minio.GetObjectOptions{})
	if err != nil {
		log.Errorf("Error retrieving s3 object %s", err.Error())
		return nil, err
	}

	byteContent, err := ioutil.ReadFile(tempFile.Name())
	if err != nil {
		log.Errorf("Error reading content from file %s", err.Error())
		return nil, err
	}

	return bytes.NewReader(byteContent), nil
}
