package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	appConfig "github.com/okoye-dev/oss-archive/internal/config"
)

type StorageInterface interface {
	UploadFile(fileName string, reader io.Reader, fileSize int64, contentType string) error
	DeleteFile(fileName string) error
	ListFiles() ([]string, error)
	GetFileSize(fileName string) (int64, error)
	GetPresignedURL(fileName string, forceDownload bool) (string, error)
}

type S3Storage struct {
	client     *s3.Client
	bucketName string
}

func NewS3Storage(cfg *appConfig.S3Config) (StorageInterface, error) {
	// Create optimized HTTP client for faster uploads
	httpClient := &http.Client{
		Timeout: 10 * time.Minute, // Long timeout for large uploads
		Transport: &http.Transport{
			MaxIdleConns:        100, 
			MaxIdleConnsPerHost: 20, 
			IdleConnTimeout:     90 * time.Second,
			DisableCompression:  false,
		},
	}

	loadOpts := []func(*config.LoadOptions) error{
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		),
		config.WithHTTPClient(httpClient), // Use optimized HTTP client
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
    // Bucket must already exist (especially for Cloudflare R2)

	return storage, nil
}

func (s *S3Storage) UploadFile(fileName string, reader io.Reader, fileSize int64, contentType string) error {
	ctx := context.Background()
	
	uploader := manager.NewUploader(s.client, func(u *manager.Uploader) {
		u.PartSize = 16 * 1024 * 1024 
		u.Concurrency = 8            
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

func (s *S3Storage) GetFileSize(fileName string) (int64, error) {
	ctx := context.Background()
	
	result, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fileName),
	})
	
	if err != nil {
		return 0, fmt.Errorf("failed to get file size: %w", err)
	}

	return aws.ToInt64(result.ContentLength), nil
}

func (s *S3Storage) GetPresignedURL(fileName string, forceDownload bool) (string, error) {
	ctx := context.Background()
	
	// Create presigned client
	presignClient := s3.NewPresignClient(s.client)
	
	// Build GetObject input
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fileName),
	}
	
	// Only add Content-Disposition header if forcing download
	if forceDownload {
		// Extract original filename for Content-Disposition header
		originalFilename := fileName
		if strings.Contains(fileName, "_") {
			// Remove UUID prefix to get original filename
			parts := strings.SplitN(fileName, "_", 2)
			if len(parts) == 2 {
				originalFilename = parts[1]
			}
		}
		input.ResponseContentDisposition = aws.String(fmt.Sprintf("attachment; filename=\"%s\"", originalFilename))
	}
	
	// Create presigned URL
	request, err := presignClient.PresignGetObject(ctx, input, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(1 * time.Hour)
	})
	
	if err != nil {
		return "", fmt.Errorf("failed to create presigned URL: %w", err)
	}
	
	return request.URL, nil
}
