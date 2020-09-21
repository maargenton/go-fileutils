package downloader_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/maargenton/fileutil/downloader"
	"github.com/maargenton/go-errors"
)

const mockError = errors.Sentinel("mockError")

var pwd, _ = os.Getwd()

type mockRequest struct {
	Method string
	Path   string
}

type mockResponse struct {
	Response     *http.Response
	ResponseBody string
}

type mockTransport struct {
	Calls          []*http.Request
	Responses      map[mockRequest]mockResponse
	RoundTripFunc  func(req *http.Request, defaultResp *http.Response) (resp *http.Response, err error)
	TransportError error
}

func (t *mockTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	t.Calls = append(t.Calls, req.Clone(req.Context()))

	var r = mockRequest{req.Method, req.URL.Path}
	if rr, ok := t.Responses[r]; ok {
		resp = &http.Response{}
		*resp = *rr.Response
		resp.Request = req
		if resp.Header == nil {
			resp.Header = make(http.Header)
		}
		resp.Header.Set("Content-Length", fmt.Sprintf("%v", len(rr.ResponseBody)))
		resp.Body = ioutil.NopCloser(strings.NewReader(rr.ResponseBody))
	} else if t.TransportError != nil {
		err = t.TransportError
		return
	}

	if t.RoundTripFunc != nil {
		resp, err = t.RoundTripFunc(req, resp)
	}

	return
}

func transportErrorClient() *downloader.Client {
	return &downloader.Client{
		HTTPClient: &http.Client{
			Transport: &mockTransport{
				TransportError: mockError,
			},
		},
	}
}
