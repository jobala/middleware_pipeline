package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/jobala/middleware_pipeline/pipeline"
)

type BearerMiddleware struct{}
type ContentTypeMiddleware struct{}

// Adds authorization header
func (s BearerMiddleware) Intercept(pipeline pipeline.Pipeline, req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", "Bearer token")

	body, _ := httputil.DumpRequest(req, true)
	log.Println(fmt.Sprintf("%s", string(body)))
	return pipeline.Next(req)
}

// Adds ContentType
func (h ContentTypeMiddleware) Intercept(pipeline pipeline.Pipeline, req *http.Request) (*http.Response, error) {
	req.Header.Add("Content-Type", "application/json")

	body, _ := httputil.DumpRequest(req, true)
	log.Println(fmt.Sprintf("%s", string(body)))
	return pipeline.Next(req)
}

func main() {
	transport := pipeline.NewCustomTransport(&BearerMiddleware{}, &ContentTypeMiddleware{})
	transport.MaxIdleConns = 10
	transport.IdleConnTimeout = 30 * time.Second

	client := &http.Client{Transport: transport}
	resp, err := client.Get("https://example.com")

	if err == nil {
		fmt.Println(resp.Status)
	}
}
