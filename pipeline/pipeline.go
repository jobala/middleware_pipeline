package pipeline

import "net/http"

type Pipeline interface {
	Next(req *http.Request) (*http.Response, error)
	Request() *http.Request
}

type Middleware interface {
	Intercept(Pipeline) (*http.Response, error)
}

type customTransport struct {
	http.Transport
	transport   http.RoundTripper
	middlewares []Middleware
}

func NewCustomTransport(middlewares ...Middleware) *customTransport {
	return &customTransport{
		middlewares: middlewares,
		transport:   http.DefaultTransport,
	}
}

type middlewarePipeline struct {
	middlewareIndex int
	transport       http.RoundTripper
	request         *http.Request
	middlewares     []Middleware
}

func (pipeline *middlewarePipeline) Request() *http.Request {
	return pipeline.request
}

func (pipeline *middlewarePipeline) incrementMiddlewareIndex() {
	pipeline.middlewareIndex++
}

func (pipeline *middlewarePipeline) Next(req *http.Request) (*http.Response, error) {
	if pipeline.middlewareIndex < len(pipeline.middlewares) {
		middleware := pipeline.middlewares[pipeline.middlewareIndex]

		pipeline.incrementMiddlewareIndex()
		return middleware.Intercept(pipeline)
	}

	return pipeline.transport.RoundTrip(req)
}

func (customTransport *customTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	pipeline := &middlewarePipeline{
		middlewareIndex: 0,
		middlewares:     customTransport.middlewares,
		transport:       customTransport.transport,
		request:         req,
	}

	reqClone := req.Clone(req.Context())
	return pipeline.Next(reqClone)
}
