// Package pipeline provides support for chaining HTTP Client middlewares
package pipeline

import "net/http"

// Pipeline interface
type Pipeline interface {
	Next(req *http.Request) (*http.Response, error)
}

type middleware interface {
	Intercept(Pipeline, *http.Request) (*http.Response, error)
}

type customTransport struct {
	http.Transport
	middlewarePipeline *middlewarePipeline
}

// MiddlewarePipeline defines the datastructure used to model the pipeline
type middlewarePipeline struct {
	middlewareIndex int
	transport       http.RoundTripper
	request         *http.Request
	middlewares     []middleware
}

// NewCustomTransport creates a transport object with a middleware pipeline
func NewCustomTransport(middlewares ...middleware) *customTransport {
	return &customTransport{
		middlewarePipeline: newMiddlewarePipeline(middlewares),
	}
}

func newMiddlewarePipeline(middlewares []middleware) *middlewarePipeline {
	return &middlewarePipeline{
		middlewareIndex: 0,
		transport:       http.DefaultTransport,
		middlewares:     middlewares,
	}
}

func (pipeline *middlewarePipeline) incrementMiddlewareIndex() {
	pipeline.middlewareIndex++
}

// Next moves the request object through middlewares in the pipeline
func (pipeline *middlewarePipeline) Next(req *http.Request) (*http.Response, error) {
	if pipeline.middlewareIndex < len(pipeline.middlewares) {
		middleware := pipeline.middlewares[pipeline.middlewareIndex]

		pipeline.incrementMiddlewareIndex()
		return middleware.Intercept(pipeline, req)
	}

	return pipeline.transport.RoundTrip(req)
}

func (transport *customTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	reqClone := req.Clone(req.Context())
	return transport.middlewarePipeline.Next(reqClone)
}
