package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	appConfig "github.com/okoye-dev/oss-archive/internal/config"
)

type StorageInterface interface {
	UploadFile(fileName string, reader io.Reader, fileSize int64, contentType string) error
	GetFile(fileName string) (io.ReadCloser, error)
	DeleteFile(fileName string) error
	ListFiles() ([]string, error)
}

type S3Storage struct {
	client     *s3.Client
	bucketName string
}

func NewS3Storage(cfg *appConfig.S3Config) (StorageInterface, error) {
	loadOpts := []func(*config.LoadOptions) error{
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		),
	}

	if cfg.Endpoint != "" {
		loadOpts = append(loadOpts, config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				if service == s3.ServiceID {
					endpoint := cfg.Endpoint
					if cfg.UseSSL {
						endpoint = "https://" + endpoint
					} else {
						endpoint = "http://" + endpoint
					}
					return aws.Endpoint{
						URL:           endpoint,
						SigningRegion: cfg.Region,
					}, nil
				}
				return aws.Endpoint{}, &aws.EndpointNotFoundError{}
			},
		)))
	}

	awsCfg, err := config.LoadDefaultConfig(context.Background(), loadOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = cfg.ForcePathStyle
	})

	storage := &S3Storage{
		client:     s3Client,
		bucketName: cfg.BucketName,
	}

	if err := storage.ensureBucket(); err != nil {
		return nil, fmt.Errorf("failed to ensure bucket exists: %w", err)
	}

	return storage, nil
}

func (s *S3Storage) ensureBucket() error {
	ctx := context.Background()
	
	_, err := s.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(s.bucketName),
	})
	
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			_, err = s.client.CreateBucket(ctx, &s3.CreateBucketInput{
				Bucket: aws.String(s.bucketName),
			})
			if err != nil {
				return fmt.Errorf("failed to create bucket: %w", err)
			}
			log.Printf("Created bucket: %s", s.bucketName)
		} else {
			return fmt.Errorf("failed to check bucket: %w", err)
		}
	}

	return nil
}

func (s *S3Storage) UploadFile(fileName string, reader io.Reader, fileSize int64, contentType string) error {
	ctx := context.Background()
	
	uploader := manager.NewUploader(s.client, func(u *manager.Uploader) {
		u.PartSize = 8 * 1024 * 1024 // 8 MiB parts
		u.Concurrency = 4            // modest parallelism
	})

	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(fileName),
		Body:        reader,
		ContentType: aws.String(contentType),
	})
	
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	log.Printf("Successfully uploaded file: %s", fileName)
	return nil
}

func (s *S3Storage) GetFile(fileName string) (io.ReadCloser, error) {
	ctx := context.Background()
	
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fileName),
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	return result.Body, nil
}

func (s *S3Storage) DeleteFile(fileName string) error {
	ctx := context.Background()
	
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fileName),
	})
	
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	log.Printf("Successfully deleted file: %s", fileName)
	return nil
}

func (s *S3Storage) ListFiles() ([]string, error) {
	ctx := context.Background()
	
	result, err := s.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucketName),
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	var files []string
	for _, obj := range result.Contents {
		files = append(files, aws.ToString(obj.Key))
	}

	return files, nil
}
