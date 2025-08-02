package http

import (
	"fmt"
	"github.com/goccy/go-json"
	"net/http"
	"time"

	"github.com/valyala/fasthttp"
)

// RequestOptions defines the options for making an HTTP request.
type RequestOptions struct {
	Method      string
	URL         string
	Body        []byte
	Headers     http.Header
	Timeout     time.Duration
	QueryParams map[string]string
}

// NewRequest makes an HTTP request based on the provided options and populates the response into the given interface.
// The 'resp' interface should be a type that can handle the response body, such as '[]byte' or 'string',
// or a custom struct where you handle the unmarshaling.
func NewRequest(requestOptions *RequestOptions, resp interface{}) error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	fresp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(fresp)

	req.SetRequestURI(requestOptions.URL)
	req.SetTimeout(requestOptions.Timeout)
	if len(requestOptions.Method) == 0 {
		requestOptions.Method = http.MethodGet
	}
	req.Header.SetMethod(requestOptions.Method)

	// Set headers
	for key, values := range requestOptions.Headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Set body
	if requestOptions.Body != nil {
		req.SetBody(requestOptions.Body)
	}

	// Set query parameters
	if len(requestOptions.QueryParams) > 0 {
		q := req.URI().QueryArgs()
		for key, value := range requestOptions.QueryParams {
			q.Add(key, value)
		}
	}

	if err := fasthttp.Do(req, fresp); err != nil {
		return fmt.Errorf("error performing request: %w", err)
	}

	statusCode := fresp.StatusCode()
	body := fresp.Body()

	if resp != nil {
		_ = json.Unmarshal(body, resp)
	}

	if statusCode >= 400 {
		return fmt.Errorf("request failed with status code: %d, response: %s", statusCode, body)
	}

	return nil
}
