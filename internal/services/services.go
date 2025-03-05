package services

import "context"

type SvcLayer struct {
	GenHTTPServer GenHTTPServer
	GenHTTPClient GenHTTPClient
}

type GenHTTPServer interface {
	Generate(ctx context.Context, serverName string) ([]byte, error)
	UploadSpec(ctx context.Context, serverName string, data []byte) error
}

type GenHTTPClient interface {
	Generate(ctx context.Context, serverName string) ([]byte, error)
}
