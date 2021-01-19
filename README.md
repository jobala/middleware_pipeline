# Middleware Pipeline

Middleware Pipeline for the Go HTTP Client.

## Getting Started

Create a middleware

```go
type AuthorizationMiddleware struct{}

func (s AuthorizationMiddleware) Intercept(pipeline pipeline.Pipeline, req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", "Bearer token")

	body, _ := httputil.DumpRequest(req, true)
	log.Println(fmt.Sprintf("%s", string(body)))

	/*
	If you want to perform an action based on the response, do the following
	
	resp, err = pipeline.Next
	// perform some action

	return resp, err
	*/
	return pipeline.Next(req)
}
```

Create a transport object

```go
// NewCustomTransport can take multiple middlewares
transport := pipeline.NewCustomTransport(&AuthorizationMiddleware{})

transport.MaxIdleConns = 10
transport.IdleConnTimeout = 30 * time.Second
```

Hook up transport with HTTP client

```go
client := &http.Client{Transport: transport}
resp, err := client.Get("https://example.com")

if err == nil {
  fmt.Println(resp.Status)
}
```
