package genhttpserver

type GenHTTPServerService struct {
}

func New() *GenHTTPServerService {
	return &GenHTTPServerService{}
}

func (s *GenHTTPServerService) Stop() error {
	return nil
}
