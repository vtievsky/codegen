// Package clienthttp provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package clienthttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/oapi-codegen/runtime"
)

const (
	BearerScopes = "bearer.Scopes"
)

// Defines values for ResponseStatusErrorCode.
const (
	Error ResponseStatusErrorCode = "error"
)

// Defines values for ResponseStatusOkCode.
const (
	Ok ResponseStatusOkCode = "ok"
)

// GenerateSpecServerHttpResponse200 defines model for GenerateSpecServerHttpResponse200.
type GenerateSpecServerHttpResponse200 struct {
	// Spec Содержимое спецификации
	Spec   []byte           `json:"spec"`
	Status ResponseStatusOk `json:"status"`
}

// GenerateSpecServerHttpResponse500 defines model for GenerateSpecServerHttpResponse500.
type GenerateSpecServerHttpResponse500 struct {
	Status ResponseStatusError `json:"status"`
}

// ResponseStatusError defines model for ResponseStatusError.
type ResponseStatusError struct {
	Code        ResponseStatusErrorCode `json:"code"`
	Description string                  `json:"description"`
}

// ResponseStatusErrorCode defines model for ResponseStatusError.Code.
type ResponseStatusErrorCode string

// ResponseStatusOk defines model for ResponseStatusOk.
type ResponseStatusOk struct {
	Code        ResponseStatusOkCode `json:"code"`
	Description string               `json:"description"`
}

// ResponseStatusOkCode defines model for ResponseStatusOk.Code.
type ResponseStatusOkCode string

// UploadSpecHttpRequest defines model for UploadSpecHttpRequest.
type UploadSpecHttpRequest struct {
	// Name Название приложения спецификация которого выгружается
	Name string `json:"name"`

	// Spec Содержимое спецификации
	Spec []byte `json:"spec"`
}

// UploadSpecHttpResponse200 defines model for UploadSpecHttpResponse200.
type UploadSpecHttpResponse200 struct {
	Status ResponseStatusOk `json:"status"`
}

// UploadSpecHttpResponse500 defines model for UploadSpecHttpResponse500.
type UploadSpecHttpResponse500 struct {
	Status ResponseStatusError `json:"status"`
}

// GenerateSpecServerHttpParams defines parameters for GenerateSpecServerHttp.
type GenerateSpecServerHttpParams struct {
	// Name Наименование сервера спецификации
	Name string `form:"name" json:"name"`
}

// UploadSpecHttpJSONRequestBody defines body for UploadSpecHttp for application/json ContentType.
type UploadSpecHttpJSONRequestBody = UploadSpecHttpRequest

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// GenerateSpecServerHttp request
	GenerateSpecServerHttp(ctx context.Context, params *GenerateSpecServerHttpParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UploadSpecHttpWithBody request with any body
	UploadSpecHttpWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UploadSpecHttp(ctx context.Context, body UploadSpecHttpJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) GenerateSpecServerHttp(ctx context.Context, params *GenerateSpecServerHttpParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGenerateSpecServerHttpRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UploadSpecHttpWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUploadSpecHttpRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UploadSpecHttp(ctx context.Context, body UploadSpecHttpJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUploadSpecHttpRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewGenerateSpecServerHttpRequest generates requests for GenerateSpecServerHttp
func NewGenerateSpecServerHttpRequest(server string, params *GenerateSpecServerHttpParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/specs/")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if queryFrag, err := runtime.StyleParamWithLocation("form", true, "name", runtime.ParamLocationQuery, params.Name); err != nil {
			return nil, err
		} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
			return nil, err
		} else {
			for k, v := range parsed {
				for _, v2 := range v {
					queryValues.Add(k, v2)
				}
			}
		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewUploadSpecHttpRequest calls the generic UploadSpecHttp builder with application/json body
func NewUploadSpecHttpRequest(server string, body UploadSpecHttpJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewUploadSpecHttpRequestWithBody(server, "application/json", bodyReader)
}

// NewUploadSpecHttpRequestWithBody generates requests for UploadSpecHttp with any type of body
func NewUploadSpecHttpRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/specs/")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// GenerateSpecServerHttpWithResponse request
	GenerateSpecServerHttpWithResponse(ctx context.Context, params *GenerateSpecServerHttpParams, reqEditors ...RequestEditorFn) (*GenerateSpecServerHttpResponse, error)

	// UploadSpecHttpWithBodyWithResponse request with any body
	UploadSpecHttpWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UploadSpecHttpResponse, error)

	UploadSpecHttpWithResponse(ctx context.Context, body UploadSpecHttpJSONRequestBody, reqEditors ...RequestEditorFn) (*UploadSpecHttpResponse, error)
}

type GenerateSpecServerHttpResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *GenerateSpecServerHttpResponse200
	JSON500      *GenerateSpecServerHttpResponse500
}

// Status returns HTTPResponse.Status
func (r GenerateSpecServerHttpResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GenerateSpecServerHttpResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UploadSpecHttpResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *UploadSpecHttpResponse200
	JSON500      *UploadSpecHttpResponse500
}

// Status returns HTTPResponse.Status
func (r UploadSpecHttpResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UploadSpecHttpResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GenerateSpecServerHttpWithResponse request returning *GenerateSpecServerHttpResponse
func (c *ClientWithResponses) GenerateSpecServerHttpWithResponse(ctx context.Context, params *GenerateSpecServerHttpParams, reqEditors ...RequestEditorFn) (*GenerateSpecServerHttpResponse, error) {
	rsp, err := c.GenerateSpecServerHttp(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGenerateSpecServerHttpResponse(rsp)
}

// UploadSpecHttpWithBodyWithResponse request with arbitrary body returning *UploadSpecHttpResponse
func (c *ClientWithResponses) UploadSpecHttpWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UploadSpecHttpResponse, error) {
	rsp, err := c.UploadSpecHttpWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUploadSpecHttpResponse(rsp)
}

func (c *ClientWithResponses) UploadSpecHttpWithResponse(ctx context.Context, body UploadSpecHttpJSONRequestBody, reqEditors ...RequestEditorFn) (*UploadSpecHttpResponse, error) {
	rsp, err := c.UploadSpecHttp(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUploadSpecHttpResponse(rsp)
}

// ParseGenerateSpecServerHttpResponse parses an HTTP response from a GenerateSpecServerHttpWithResponse call
func ParseGenerateSpecServerHttpResponse(rsp *http.Response) (*GenerateSpecServerHttpResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GenerateSpecServerHttpResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest GenerateSpecServerHttpResponse200
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest GenerateSpecServerHttpResponse500
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseUploadSpecHttpResponse parses an HTTP response from a UploadSpecHttpWithResponse call
func ParseUploadSpecHttpResponse(rsp *http.Response) (*UploadSpecHttpResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &UploadSpecHttpResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest UploadSpecHttpResponse200
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest UploadSpecHttpResponse500
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}
