package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jobala/middleware_pipeline/pipeline"
)

type StarsMiddleware struct{}
type HashMiddleware struct{}

func (s StarsMiddleware) Intercept(pipeline pipeline.Pipeline) (*http.Response, error) {
	fmt.Println("*******************")

	req := pipeline.Request()

	return pipeline.Next(req)
}

func (h HashMiddleware) Intercept(pipeline pipeline.Pipeline) (*http.Response, error) {
	fmt.Println("####################")

	req := pipeline.Request()

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
