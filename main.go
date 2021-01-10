package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/jobala/middleware_pipeline/pipeline"
)

type StarsMiddleware struct{}
type HashMiddleware struct{}

func (s StarsMiddleware) Intercept(pipeline pipeline.Pipeline) (*http.Response, error) {
	req := pipeline.Request()
	req.Header.Add("Authorization", "Bearer token")

	body, _ := httputil.DumpRequest(req, true)
	log.Println(fmt.Sprintf("%s", string(body)))
	return pipeline.Next(req)
}

func (h HashMiddleware) Intercept(pipeline pipeline.Pipeline) (*http.Response, error) {
	req := pipeline.Request()
	req.Header.Add("Content-Type", "application/json")

	body, _ := httputil.DumpRequest(req, true)
	log.Println(fmt.Sprintf("%s", string(body)))
	return pipeline.Next(req)
}

func main() {
	transport := pipeline.NewCustomTransport(&StarsMiddleware{}, &HashMiddleware{})
	transport.MaxIdleConns = 10
	transport.IdleConnTimeout = 30 * time.Second

	client := &http.Client{Transport: transport}
	resp, err := client.Get("https://example.com")

	if err == nil {
		fmt.Println(resp.Status)
	}
}
