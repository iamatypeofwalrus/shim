[![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoib0UxQzQ0NXV5R0NhUUpMZHV0dmU4WS9yNG43Ynlvb3Y2WmdNZUhQYWEzMXdkMjJCNFgvUkJvSlY5aEZ6R0wyUi9Ud1B1Vll0R2FIQytpdGU3QllDUFE0PSIsIml2UGFyYW1ldGVyU3BlYyI6IlNDdzM1NmY0ZU5SWjV2aE4iLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=master)](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoib0UxQzQ0NXV5R0NhUUpMZHV0dmU4WS9yNG43Ynlvb3Y2WmdNZUhQYWEzMXdkMjJCNFgvUkJvSlY5aEZ6R0wyUi9Ud1B1Vll0R2FIQytpdGU3QllDUFE0PSIsIml2UGFyYW1ldGVyU3BlYyI6IlNDdzM1NmY0ZU5SWjV2aE4iLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=master)
[![GoDoc](https://godoc.org/github.com/iamatypeofwalrus/shim?status.svg)](https://godoc.org/github.com/iamatypeofwalrus/shim)
[![Go Report Card](https://goreportcard.com/badge/github.com/iamatypeofwalrus/shim)](https://goreportcard.com/report/github.com/iamatypeofwalrus/shim)

# Shim
Shim is a thin layer between API Gateway integration requests via Lambda and the standard library `http.Handler` interface. It allows you to write plain ol' Go and run it on Lambda with minimal modifications. Bring your own router!

Shim uses [Go modules](https://github.com/golang/go/wiki/Modules) to model its dependencies.

## Example
For an extensive example on how `shim` fits in with other AWS serverless tooling like [SAM Local](https://github.com/awslabs/aws-sam-local) and the [Serverless Application Model (SAM) specification](https://github.com/awslabs/serverless-application-model) head over to the [this example in the wiki](https://github.com/iamatypeofwalrus/shim/wiki/Example:-AWS-Sam-Local)

### Note: API Gateway
Make sure that [proxy pass integration in API Gateway](https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-set-up-simple-proxy.html) is enabled to make sure your application receives every request sent to your API Gateway endpoint.

### Code
#### With Rest API (API Gateway Proxy Request, Response)
```go
package main

import (
  "fmt"
  "net/http"

  "github.com/aws/aws-lambda-go/lambda"

  "github.com/iamatypeofwalrus/shim"
)

func main() {
  // Create a router as normal. Any router that satisfies the http.Handler interface
  // is accepted!
  mux := http.NewServeMux()
  mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
    fmt.Fprint(w, "hello, world")
  })

  s := shim.HandleRestApiRequests(mux)

  // Pass your router to shim and let Lambda handle the rest
  lambda.Start(s.Handle)
}
```

#### With HTTP API (API Gateway V2 Request, Response)
```go
package main

import (
  "fmt"
  "net/http"

  "github.com/aws/aws-lambda-go/lambda"

  "github.com/iamatypeofwalrus/shim"
)

func main() {
  // Create a router as normal. Any router that satisfies the http.Handler interface
  // is accepted!
  mux := http.NewServeMux()
  mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
    fmt.Fprint(w, "hello, world")
  })

  s := shim.HandleHttpApiRequests(mux)

  // Pass your router to shim and let Lambda handle the rest
  lambda.Start(s.Handle)
}
```

### With Debugging Logger
You can pull logs from various steps in the shim by passing the `SetDebugLogger` option. [It accepts any logger that provides `Printf`](https://github.com/iamatypeofwalrus/shim/blob/56bb8c10bbb8e36d964551ceace772f675141ec8/log.go#L5) functions a lá the standard library logger.

```go
func main() {
  ...

  l := log.New(os.Stdout, "", log.LstdFlags)
  shim := shim.New(
    nil, // or your mux
    shim.SetDebugLogger(l)
  )

  ...
}
```

There is also a shim for the `slog` package

```go
func main() {
  ...

  // Debug with Slog calls the Debug method on slog. Slog defaults to INFO, so it needs to be set to Debug so you can see the messages
  logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
  shim := shim.New(
    nil, // or your mux
    shim.SetDebugWithSlog(l)
  )

  ...
}
```
