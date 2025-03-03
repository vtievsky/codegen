codegen-server:
	@oapi-codegen --config=./testdata/openapi/oapi-codegen-server.yaml ./docs/openapi/swagger.yaml

codegen-client:
	@oapi-codegen --config=./testdata/openapi/oapi-codegen-client.yaml ./docs/openapi/swagger.yaml	