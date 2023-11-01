// This file was auto-generated by Fern from our API Definition.

package documents

import (
	bytes "bytes"
	context "context"
	json "encoding/json"
	errors "errors"
	fmt "fmt"
	vellumclientgo "github.com/vellum-ai/vellum-client-go"
	core "github.com/vellum-ai/vellum-client-go/core"
	io "io"
	multipart "mime/multipart"
	http "net/http"
	url "net/url"
)

type Client struct {
	baseURL string
	caller  *core.Caller
	header  http.Header
}

func NewClient(opts ...core.ClientOption) *Client {
	options := core.NewClientOptions()
	for _, opt := range opts {
		opt(options)
	}
	return &Client{
		baseURL: options.BaseURL,
		caller:  core.NewCaller(options.HTTPClient),
		header:  options.ToHeader(),
	}
}

// Used to list documents. Optionally filter on supported fields.
func (c *Client) List(ctx context.Context, request *vellumclientgo.DocumentsListRequest) (*vellumclientgo.PaginatedSlimDocumentList, error) {
	baseURL := "https://api.vellum.ai"
	if c.baseURL != "" {
		baseURL = c.baseURL
	}
	endpointURL := baseURL + "/" + "v1/documents"

	queryParams := make(url.Values)
	if request.DocumentIndexId != nil {
		queryParams.Add("document_index_id", fmt.Sprintf("%v", *request.DocumentIndexId))
	}
	if request.Limit != nil {
		queryParams.Add("limit", fmt.Sprintf("%v", *request.Limit))
	}
	if request.Offset != nil {
		queryParams.Add("offset", fmt.Sprintf("%v", *request.Offset))
	}
	if request.Ordering != nil {
		queryParams.Add("ordering", fmt.Sprintf("%v", *request.Ordering))
	}
	if len(queryParams) > 0 {
		endpointURL += "?" + queryParams.Encode()
	}

	var response *vellumclientgo.PaginatedSlimDocumentList
	if err := c.caller.Call(
		ctx,
		&core.CallParams{
			URL:      endpointURL,
			Method:   http.MethodGet,
			Headers:  c.header,
			Response: &response,
		},
	); err != nil {
		return nil, err
	}
	return response, nil
}

// Update a Document, keying off of its Vellum-generated ID. Particularly useful for updating its metadata.
//
// A UUID string identifying this document.
func (c *Client) PartialUpdate(ctx context.Context, id string, request *vellumclientgo.PatchedDocumentUpdateRequest) (*vellumclientgo.DocumentRead, error) {
	baseURL := "https://api.vellum.ai"
	if c.baseURL != "" {
		baseURL = c.baseURL
	}
	endpointURL := fmt.Sprintf(baseURL+"/"+"v1/documents/%v", id)

	var response *vellumclientgo.DocumentRead
	if err := c.caller.Call(
		ctx,
		&core.CallParams{
			URL:      endpointURL,
			Method:   http.MethodPatch,
			Headers:  c.header,
			Request:  request,
			Response: &response,
		},
	); err != nil {
		return nil, err
	}
	return response, nil
}

// A UUID string identifying this document.
func (c *Client) Destroy(ctx context.Context, id string) error {
	baseURL := "https://api.vellum.ai"
	if c.baseURL != "" {
		baseURL = c.baseURL
	}
	endpointURL := fmt.Sprintf(baseURL+"/"+"v1/documents/%v", id)

	if err := c.caller.Call(
		ctx,
		&core.CallParams{
			URL:     endpointURL,
			Method:  http.MethodDelete,
			Headers: c.header,
		},
	); err != nil {
		return err
	}
	return nil
}

// Upload a document to be indexed and used for search.
//
// **Note:** Uses a base url of `https://documents.vellum.ai`.
func (c *Client) Upload(ctx context.Context, contents io.Reader, request *vellumclientgo.UploadDocumentBodyRequest) (*vellumclientgo.UploadDocumentResponse, error) {
	baseURL := "https://documents.vellum.ai"
	if c.baseURL != "" {
		baseURL = c.baseURL
	}
	endpointURL := baseURL + "/" + "v1/upload-document"

	errorDecoder := func(statusCode int, body io.Reader) error {
		raw, err := io.ReadAll(body)
		if err != nil {
			return err
		}
		apiError := core.NewAPIError(statusCode, errors.New(string(raw)))
		decoder := json.NewDecoder(bytes.NewReader(raw))
		switch statusCode {
		case 400:
			value := new(vellumclientgo.BadRequestError)
			value.APIError = apiError
			if err := decoder.Decode(value); err != nil {
				return apiError
			}
			return value
		case 404:
			value := new(vellumclientgo.NotFoundError)
			value.APIError = apiError
			if err := decoder.Decode(value); err != nil {
				return apiError
			}
			return value
		case 500:
			value := new(vellumclientgo.InternalServerError)
			value.APIError = apiError
			if err := decoder.Decode(value); err != nil {
				return apiError
			}
			return value
		}
		return apiError
	}

	var response *vellumclientgo.UploadDocumentResponse
	requestBuffer := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(requestBuffer)
	contentsFilename := "contents_filename"
	if named, ok := contents.(interface{ Name() string }); ok {
		contentsFilename = named.Name()
	}
	contentsPart, err := writer.CreateFormFile("contents", contentsFilename)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(contentsPart, contents); err != nil {
		return nil, err
	}
	if request.AddToIndexNames != nil {
		if err := core.WriteMultipartJSON(writer, "add_to_index_names", request.AddToIndexNames); err != nil {
			return nil, err
		}
	}
	if request.ExternalId != nil {
		if err := writer.WriteField("external_id", fmt.Sprintf("%v", *request.ExternalId)); err != nil {
			return nil, err
		}
	}
	if err := writer.WriteField("label", fmt.Sprintf("%v", request.Label)); err != nil {
		return nil, err
	}
	if request.Keywords != nil {
		if err := core.WriteMultipartJSON(writer, "keywords", request.Keywords); err != nil {
			return nil, err
		}
	}
	if request.Metadata != nil {
		if err := writer.WriteField("metadata", fmt.Sprintf("%v", *request.Metadata)); err != nil {
			return nil, err
		}
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	c.header.Set("Content-Type", writer.FormDataContentType())

	if err := c.caller.Call(
		ctx,
		&core.CallParams{
			URL:          endpointURL,
			Method:       http.MethodPost,
			Headers:      c.header,
			Request:      requestBuffer,
			Response:     &response,
			ErrorDecoder: errorDecoder,
		},
	); err != nil {
		return nil, err
	}
	return response, nil
}