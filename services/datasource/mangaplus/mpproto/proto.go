package mpproto

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	proto "google.golang.org/protobuf/proto"
)

// https://github.com/protocolbuffers/protobuf/releases
// google.golang.org/protobuf/cmd/protoc-gen-go

//go:generate protoc --go_out=.. mpproto.proto

type Client struct {
	sessionToken string
	httpClient   http.Client
}

func NewClient(c *http.Client) *Client {
	return &Client{
		sessionToken: uuid.NewString(),
		httpClient:   *c,
	}
}

func (c *Client) Get(path string, a ...interface{}) (*SuccessResult, error) {
	return c.request("GET", fmt.Sprintf(path, a...), nil)
}

func (c *Client) Put(body url.Values, path string, a ...interface{}) (*SuccessResult, error) {
	return c.request("PUT", fmt.Sprintf(path, a...), nil)
}

func (c *Client) request(method, url string, body url.Values) (*SuccessResult, error) {
	req, err := http.NewRequest(method, url, getBody(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("SESSION-TOKEN", c.sessionToken)

	r, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	resp := &Response{}

	err = proto.Unmarshal(b, resp)
	if err != nil {
		return nil, err
	}

	if resp.GetSuccess() == nil {
		return nil, resp.GetError()
	}

	return resp.GetSuccess(), nil
}

func (e *ErrorResult) Error() string {
	return e.String()
}

func getBody(v url.Values) io.Reader {
	if v == nil {
		return http.NoBody
	}
	return bytes.NewBufferString(v.Encode())
}
