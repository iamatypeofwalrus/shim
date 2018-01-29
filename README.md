[![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoibkxJRnI0U3VBaTBEWnRpQlBLOGNMOE1Lb2dTSzlUTDVJRHJTcHdmOTk1ZFVKSGVwbzFjdm5ybG9WZTZFUWtpaFdoSnh0RVNROW9aTVFhZzVIb1BOVHpNPSIsIml2UGFyYW1ldGVyU3BlYyI6Ikk2VFlCMEh3M3kzRDJuQnQiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=master)](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoibkxJRnI0U3VBaTBEWnRpQlBLOGNMOE1Lb2dTSzlUTDVJRHJTcHdmOTk1ZFVKSGVwbzFjdm5ybG9WZTZFUWtpaFdoSnh0RVNROW9aTVFhZzVIb1BOVHpNPSIsIml2UGFyYW1ldGVyU3BlYyI6Ikk2VFlCMEh3M3kzRDJuQnQiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=master)
[![GoDoc](https://godoc.org/github.com/iamatypeofwalrus/shim?status.svg)](https://godoc.org/github.com/iamatypeofwalrus/shim)
[![Go Report Card](https://goreportcard.com/badge/github.com/iamatypeofwalrus/shim)](https://goreportcard.com/report/github.com/iamatypeofwalrus/shim)

# Shim
Shim is a thin layer between API Gateway integration requests via Lambda and the standard library `http.Handler` interface. It allows you to write plain ol' Go and run it on Lambda with minimal modifications. Bring your own router!

Shim uses [`dep`](https://golang.github.io/dep/) to manage its dependencies. You can add `shim` to your dep project by running:

```sh
dep ensure -add github.com/iamatypeofwalrus/shim
```

## Usage
### CloudFormation
You'll want to use the [proxy pass integration](https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-set-up-simple-proxy.html) with API Gateway to make sure your application receives every request sent to your API Gateway endpoint.

```
# Here we're using the SAM specification to define our function
#
# Note: You need both the Root AND the Greedy event in order to capture all
#       events sent to your web app.
YourFunction:
  Type: AWS::Serverless::Function
  Properties:
    Handler: main
    Runtime: go1.x
    Role: ...
    Events:
      ProxyApiRoot:
        Type: Api
        Properties:
          Path: /
          Method: ANY
      ProxyApiGreedy:
        Type: Api
        Properties:
          Path: /{proxy+}
          Method: ANY
```
### Go Code
#### With your own router
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

  s := shim.New(mux)

  // Pass your router to shim and let Lambda handle the rest
  lambda.Start(s.Handle)
}
```

#### With the default router
```go
package main

import (
  "fmt"
  "net/http"

  "github.com/aws/aws-lambda-go/lambda"

  "github.com/iamatypeofwalrus/shim"
)

func main() {
  // Shim works with the http.DefaultServeMux. Create routes and handlers against the router normal.
  http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
    fmt.Fprint(w, "hello, world")
  })

  // Simply pass nil to shim to use the http.DefaultServeMux
  s := shim.New(nil)

  lambda.Start(s.Handle)
}
```

#### With Debugging Logger
You can pull logs from various steps in the shim by passing the `SetDebugLogger` option. [It accepts any logger that provides
the `Println` and `Printf`](https://github.com/iamatypeofwalrus/shim/blob/56bb8c10bbb8e36d964551ceace772f675141ec8/log.go#L5) functions a lá the standard library logger.

```go
func main() {
  ...

  l := log.New(os.Stdout, "", log.LstdFlags)
  shim := shim.New(
    nil, // or your mux
    shim.SetDebugLogger(l)
  )

  lambda.Start(shim.Handle)
}
```
