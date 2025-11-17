package riot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

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

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
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

type RiotClientInterface interface {
	GetAccountByRiotId(ctx context.Context, params AccountByRiotIdRequestParams) (AccountByRiotIdResponse, error)

	GetQueueEntriesByPlayerUuid(ctx context.Context, params QueueEntriesByPlayerUuidParams) ([]*QueueResponse, error)
}

func (c *Client) GetAccountByRiotId(
	ctx context.Context,
	params AccountByRiotIdRequestParams,
) (*AccountByRiotIdResponse, error) {
	var err error

	serverURL, err := url.Parse(c.Server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("riot/account/v1/accounts/by-riot-id/%s/%s", params.Name, params.Tagline)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, c.RequestEditors); err != nil {
		return nil, err
	}

	response, err :=  c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := io.ReadAll(response.Body)
	defer func() { _ = response.Body.Close() }()
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		var dest Status
		err = json.Unmarshal(bodyBytes, &dest)
		if err != nil {
			return nil, err
		}

		return nil, errors.New(fmt.Sprintf("HTTP StatusCode: %d - %s", dest.StatusCode, dest.Message))
	}

	var dest AccountByRiotIdResponse
	err = json.Unmarshal(bodyBytes, &dest)

	if err != nil {
	  return nil, err
	}

	return &dest, nil

}
