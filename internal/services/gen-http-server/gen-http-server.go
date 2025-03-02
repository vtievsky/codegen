package genhttpserver

import "go.uber.org/zap"

type GenHTTPServerServiceOpts struct {
	Logger *zap.Logger
}

type GenHTTPServerService struct {
	logger *zap.Logger
}

func New(opts *GenHTTPServerServiceOpts) *GenHTTPServerService {
	return &GenHTTPServerService{
		logger: opts.Logger,
	}
}

func (s *GenHTTPServerService) Start() error {
	s.logger.Info("codegen start")

	return nil
}

func (s *GenHTTPServerService) Stop() error {
	s.logger.Info("codegen stop")

	return nil
}
