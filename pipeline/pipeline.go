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

type middlewareChain struct {
	middlewareIndex int
	transport       http.RoundTripper
	request         *http.Request
	middlewares     []Middleware
}

func (chain *middlewareChain) Request() *http.Request {
	return chain.request
}

func (chain *middlewareChain) Next(req *http.Request) (*http.Response, error) {
	if chain.middlewareIndex < len(chain.middlewares) {
		c := &middlewareChain{
			middlewareIndex: chain.middlewareIndex + 1,
			middlewares:     chain.middlewares,
			transport:       chain.transport,
			request:         req,
		}

		middleware := chain.middlewares[chain.middlewareIndex]
		return middleware.Intercept(c)
	}

	return chain.transport.RoundTrip(req)
}

func (customTransport customTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	chain := &middlewareChain{
		middlewareIndex: 0,
		middlewares:     customTransport.middlewares,
		transport:       customTransport.transport,
		request:         req,
	}

	reqClone := req.Clone(req.Context())
	return chain.Next(reqClone)
}
