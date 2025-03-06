package specstorage

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	BuffSize  int    = 1024
	awsRegion string = "ap-southeast-1"
)

type SpecStoreClientOpts struct {
	AccessKey string
	SecretKey string
	URL       string
}

type SpecStoreClient struct {
	client *s3.Client
}

func New(opts *SpecStoreClientOpts) *SpecStoreClient {
	client := s3.NewFromConfig(aws.Config{ //nolint:exhaustruct
		Region:      awsRegion,
		Credentials: credentials.NewStaticCredentialsProvider(opts.AccessKey, opts.SecretKey, ""),
		EndpointResolverWithOptions: aws.EndpointResolverWithOptionsFunc( //nolint:staticcheck
			func(service string, region string, options ...any) (aws.Endpoint, error) { //nolint:staticcheck
				return aws.Endpoint{ //nolint:exhaustruct,staticcheck
					PartitionID:       "aws",
					SigningRegion:     awsRegion,
					URL:               opts.URL,
					HostnameImmutable: true,
				}, nil
			},
		),
	})

	return &SpecStoreClient{
		client: client,
	}
}

func (s *SpecStoreClient) Upload(ctx context.Context, specBucket, specName string, data []byte) error {
	uploader := manager.NewUploader(s.client)

	if _, err := uploader.Upload(ctx, &s3.PutObjectInput{ //nolint:exhaustruct
		Bucket: aws.String(specBucket),
		Key:    aws.String(specName),
		Body:   bytes.NewReader(data),
	}); err != nil {
		return fmt.Errorf("failed to upload data | %w", err)
	}

	return nil
}

func (s *SpecStoreClient) Download(ctx context.Context, specBucket, specName string) ([]byte, error) {
	buf := make([]byte, BuffSize)
	specWriter := manager.NewWriteAtBuffer(buf)
	downloader := manager.NewDownloader(s.client)

	if _, err := downloader.Download(ctx, specWriter, &s3.GetObjectInput{ //nolint:exhaustruct
		Bucket: aws.String(specBucket),
		Key:    aws.String(specName),
	}); err != nil {
		return nil, fmt.Errorf("failed to download data | %w", err)
	}

	return specWriter.Bytes(), nil
}
