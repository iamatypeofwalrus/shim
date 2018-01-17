# Shim
Shim is a thin layer between the Lambda API Gateway integrations and the standard library `http.Handler` interface. It allows you to write plain ol' Go and run it on Lambda with minimal modifications.

Bring your own router.

## Usage
### Cloudformation
You'll want to use the [proxy pass integration])(https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-set-up-simple-proxy.html) with API Gateway to make sure you application receives every request.

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
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "hello, world")
	})

	shim := shim.New(mux)
	lambda.Start(shim.Handle)
}
```