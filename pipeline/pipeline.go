package pipeline

import "net/http"

type Pipeline interface {
	Next(req *http.Request) (*http.Response, error)
}

type Middleware interface {
	Intercept(Pipeline, *http.Request) (*http.Response, error)
}

type customTransport struct {
	http.Transport
	middlewarePipeline *MiddlewarePipeline
}

type MiddlewarePipeline struct {
	middlewareIndex int
	transport       http.RoundTripper
	request         *http.Request
	middlewares     []Middleware
}

func NewCustomTransport(middlewares ...Middleware) *customTransport {
	return &customTransport{
		middlewarePipeline: &MiddlewarePipeline{
			middlewareIndex: 0,
			transport:       http.DefaultTransport,
			middlewares:     middlewares,
		},
	}
}

func (pipeline *MiddlewarePipeline) incrementMiddlewareIndex() {
	pipeline.middlewareIndex++
}

func (pipeline *MiddlewarePipeline) Next(req *http.Request) (*http.Response, error) {
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
