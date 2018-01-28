[![GoDoc](https://godoc.org/github.com/iamatypeofwalrus/shim?status.svg)](https://godoc.org/github.com/iamatypeofwalrus/shim)
![AWS CodeBuild Status](https://s3-us-west-2.amazonaws.com/codefactory-us-west-2-prod-default-build-badges/passing.svg)

# Shim
Shim is a thin layer between API Gateway integration requests via Lambda and the standard library `http.Handler` interface. It allows you to write plain ol' Go and run it on Lambda with minimal modifications. Bring your own router!

## Usage
### Cloudformation
You'll want to use the [proxy pass integration](https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-set-up-simple-proxy.html) with API Gateway to make sure your application receives every request sent to API Gateway endpoint.

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

  // Pass your router to shim and let Lambda handle the rest
  lambda.Start(shim.New(mux).Handle)
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
    mux,
    shim.SetDebugLogger(l)
  )
  lambda.Start(shim.Handle)
}
```
