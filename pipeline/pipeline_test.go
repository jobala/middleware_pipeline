package pipeline_test

import (
	"net/http"
	"testing"

	"github.com/jobala/middleware_pipeline/pipeline"
)

type TestMiddleware struct{}

func (middleware TestMiddleware) Intercept(pipeline pipeline.Pipeline, req *http.Request) (*http.Response, error) {
	req.Header.Add("test", "test-header")

	return pipeline.Next(req)
}

func TestCanInterceptRequests(t *testing.T) {
	transport := pipeline.NewCustomTransport(&TestMiddleware{})
	client := &http.Client{Transport: transport}
	resp, _ := client.Get("https://example.com")

	expect := "test-header"
	got := resp.Request.Header.Get("test")

	if expect != got {
		t.Errorf("Expected: %v, but received: %v", expect, got)
	}
}
