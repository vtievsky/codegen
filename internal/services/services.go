package services

import "context"

type SvcLayer struct {
	GenHTTPServer GenHTTPServer
}

type GenHTTPServer interface {
	Generate(ctx context.Context, serverName string) ([]byte, error)
	UploadSpec(ctx context.Context, serverName string, data []byte) error
}
