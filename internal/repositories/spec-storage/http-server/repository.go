package specstoragehttpserver

import "context"

const (
	PackageName = "serverhttp"
)

type Client interface {
	Upload(ctx context.Context, specbucket, specname string, data []byte) error
	Download(ctx context.Context, specBucket, specName string) ([]byte, error)
}

type SpecHTTPServerOpts struct {
	Client Client
}

type SpecHTTPServer struct {
	client Client
}

func New(opts *SpecHTTPServerOpts) *SpecHTTPServer {
	return &SpecHTTPServer{
		client: opts.Client,
	}
}

func (s *SpecHTTPServer) Upload(ctx context.Context, serverName string, data []byte) error {
	return s.client.Upload(ctx, PackageName, serverName, data)
}

func (s *SpecHTTPServer) Download(ctx context.Context, serverName string) ([]byte, error) {
	return s.client.Download(ctx, PackageName, serverName)
}
