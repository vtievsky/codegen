package genhttpserver

import (
	"context"
	"fmt"
	"go/format"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
	"go.uber.org/zap"
)

var (
	generateOptions = codegen.GenerateOptions{
		IrisServer:    false,
		ChiServer:     false,
		FiberServer:   false,
		EchoServer:    true,
		GinServer:     false,
		GorillaServer: false,
		StdHTTPServer: false,
		Strict:        true,
		Client:        false,
		Models:        true,
		EmbeddedSpec:  true,
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

const (
	packageName = "serverhttp"
)

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

func (s *GenHTTPServerService) GenHTTPServer(ctx context.Context, data []byte) ([]byte, error) {
	swagger, err := s.swaggerFromData(data)
	if err != nil {
		return nil, err
	}

	transportCode, err := codegen.Generate(swagger, codegen.Configuration{
		PackageName:          packageName,
		Generate:             generateOptions,
		Compatibility:        compatibilityOptions,
		OutputOptions:        outputOptions,
		ImportMapping:        nil,
		AdditionalImports:    nil,
		NoVCSVersionOverride: nil,
	})
	if err != nil {
		return nil, err
	}

	result, err := s.formatSource(transportCode)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *GenHTTPServerService) swaggerFromData(data []byte) (*openapi3.T, error) {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	swagger, err := loader.LoadFromData(data)
	if err != nil {
		return nil, err
	}

	return swagger, nil
}

func (s *GenHTTPServerService) formatSource(src string) ([]byte, error) {
	result, err := format.Source([]byte(src))
	if err != nil {
		return nil, fmt.Errorf("failed to source code formatting | %w", err)
	}

	// result, err = imports.Process(filename, result, nil)
	// if err != nil {
	// 	return nil, fmt.Errorf("imports.Process error | %w", err)
	// }

	return result, nil
}
