package genhttpclient

import (
	"context"
	"fmt"
	"go/format"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
	"go.uber.org/zap"
)

var (
	generateOptions = codegen.GenerateOptions{
		IrisServer:    false,
		ChiServer:     false,
		FiberServer:   false,
		EchoServer:    false,
		GinServer:     false,
		GorillaServer: false,
		StdHTTPServer: false,
		Strict:        false,
		Client:        true,
		Models:        true,
		EmbeddedSpec:  false,
	}
	compatibilityOptions = codegen.CompatibilityOptions{
		OldMergeSchemas:                    false,
		OldEnumConflicts:                   false,
		OldAliasing:                        false,
		DisableFlattenAdditionalProperties: false,
		DisableRequiredReadOnlyAsPointer:   false,
		AlwaysPrefixEnumValues:             false,
		ApplyChiMiddlewareFirstToLast:      false,
		ApplyGorillaMiddlewareFirstToLast:  false,
		CircularReferenceLimit:             0,
		AllowUnexportedStructFieldNames:    false,
	}
	outputOptionsOverlay = codegen.OutputOptionsOverlay{
		Path:   "",
		Strict: nil,
	}
	outputOptions = codegen.OutputOptions{
		SkipFmt:                   false,
		SkipPrune:                 false,
		IncludeTags:               nil,
		ExcludeTags:               nil,
		IncludeOperationIDs:       nil,
		ExcludeOperationIDs:       nil,
		UserTemplates:             nil, // TODO
		ExcludeSchemas:            nil,
		ResponseTypeSuffix:        "",
		ClientTypeName:            "",
		InitialismOverrides:       false,
		NullableType:              false,
		DisableTypeAliasesForType: nil,
		NameNormalizer:            "",
		Overlay:                   outputOptionsOverlay,
	}
)

type SpecStorage interface {
	Upload(ctx context.Context, specname string, data []byte) error
	Download(ctx context.Context, specName string) ([]byte, error)
}

type GenHTTPClientServiceOpts struct {
	Logger      *zap.Logger
	SpecStorage SpecStorage
}

type GenHTTPClientService struct {
	logger      *zap.Logger
	specStorage SpecStorage
}

func New(opts *GenHTTPClientServiceOpts) *GenHTTPClientService {
	return &GenHTTPClientService{
		logger:      opts.Logger,
		specStorage: opts.SpecStorage,
	}
}

func (s *GenHTTPClientService) Generate(ctx context.Context, serverName string) ([]byte, error) {
	const op = "GenHTTPClientService.Generate"

	data, err := s.loadSpec(ctx, serverName)
	if err != nil {
		return nil, fmt.Errorf("failed to load spec | %s:%w", op, err)
	}

	swagger, err := s.swaggerFromData(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse spec | %s:%w", op, err)
	}

	packageName := strings.ReplaceAll(serverName, "-", "")

	transportCode, err := codegen.Generate(swagger, codegen.Configuration{
		PackageName:          fmt.Sprintf("%shttpclient", packageName),
		Generate:             generateOptions,
		Compatibility:        compatibilityOptions,
		OutputOptions:        outputOptions,
		ImportMapping:        nil,
		AdditionalImports:    nil,
		NoVCSVersionOverride: nil,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate transport code | %s:%w", op, err)
	}

	resp, err := s.formatSource(transportCode)
	if err != nil {
		return nil, fmt.Errorf("failed to format transport code | %s:%w", op, err)
	}

	return resp, nil
}

func (s *GenHTTPClientService) swaggerFromData(data []byte) (*openapi3.T, error) {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	swagger, err := loader.LoadFromData(data)
	if err != nil {
		return nil, err
	}

	return swagger, nil
}

func (s *GenHTTPClientService) formatSource(src string) ([]byte, error) {
	resp, err := format.Source([]byte(src))
	if err != nil {
		return nil, err
	}

	// result, err = imports.Process(filename, result, nil)
	// if err != nil {
	// 	return nil, fmt.Errorf("imports.Process error | %w", err)
	// }

	return resp, nil
}

func (s *GenHTTPClientService) loadSpec(ctx context.Context, serverName string) ([]byte, error) {
	return s.specStorage.Download(ctx, serverName)
}
