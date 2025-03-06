package httptransport

import (
	"context"

	"github.com/vtievsky/codegen-svc/gen/serverhttp"
	"github.com/vtievsky/codegen-svc/internal/services"
)

type Transport struct {
	services *services.SvcLayer
}

func New(services *services.SvcLayer) *Transport {
	return &Transport{
		services: services,
	}
}

func (s *Transport) GenerateSpecServerHttp(
	ctx context.Context,
	request serverhttp.GenerateSpecServerHttpRequestObject,
) (serverhttp.GenerateSpecServerHttpResponseObject, error) {
	data, err := s.services.GenHTTPServer.Generate(ctx, request.Name)
	if err != nil {
		return serverhttp.GenerateSpecServerHttp500JSONResponse{ //nolint:nilerr
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		}, nil
	}

	return serverhttp.GenerateSpecServerHttp200JSONResponse{
		Spec: data,
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	}, nil
}

func (s *Transport) GenerateSpecClientHttp(
	ctx context.Context,
	request serverhttp.GenerateSpecClientHttpRequestObject,
) (serverhttp.GenerateSpecClientHttpResponseObject, error) {
	data, err := s.services.GenHTTPClient.Generate(ctx, request.Name)
	if err != nil {
		return serverhttp.GenerateSpecClientHttp500JSONResponse{ //nolint:nilerr
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		}, nil
	}

	return serverhttp.GenerateSpecClientHttp200JSONResponse{
		Spec: data,
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	}, nil
}

func (s *Transport) UploadSpecHttp(
	ctx context.Context,
	request serverhttp.UploadSpecHttpRequestObject,
) (serverhttp.UploadSpecHttpResponseObject, error) {
	err := s.services.GenHTTPServer.UploadSpec(ctx, request.Body.Name, request.Body.Spec)
	if err != nil {
		return serverhttp.UploadSpecHttp500JSONResponse{ //nolint:nilerr
			Status: serverhttp.ResponseStatusError{
				Code:        serverhttp.Error,
				Description: err.Error(),
			},
		}, nil
	}

	return serverhttp.UploadSpecHttp200JSONResponse{
		Status: serverhttp.ResponseStatusOk{
			Code:        serverhttp.Ok,
			Description: "",
		},
	}, nil
}
