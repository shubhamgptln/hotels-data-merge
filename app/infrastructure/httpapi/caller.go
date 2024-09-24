package httpapi

import (
	"context"
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"
)

// ErrResponseNotOk ...
var ErrResponseNotOk = errors.New("response status is not ok")

type (
	// API ...
	API interface {
		BuildRequest(ctx context.Context) (*http.Request, error)
		ParseResponse(ctx context.Context, req *http.Request, resp *http.Response) error
	}

	// Caller ...
	Caller interface {
		Call(ctx context.Context, api API) (successful bool, err error)
	}
)

// NewCaller ...
func NewCaller(c *http.Client) Caller {
	return &caller{c}
}

// Caller ...
type caller struct {
	httpClient *http.Client
}

// Call ...
func (c caller) Call(ctx context.Context, api API) (successful bool, err error) {
	req, err := api.BuildRequest(ctx)
	if err != nil {
		return false, err
	}

	req = req.WithContext(ctx)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logrus.Errorf("http err :%v", err.Error())
		return false, err
	}
	if err := api.ParseResponse(ctx, req, resp); err != nil {
		logrus.Errorf("http resp parse err :%v", err.Error())
		return false, err
	}
	return true, nil
}
