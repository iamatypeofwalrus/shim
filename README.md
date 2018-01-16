# Shim
Shim is a thin layer between the Lambda API Gateway integrations and the standard library `http.Handler` interface. It allows you to write plain ol' Go and run it on Lambda with minimal modifications.

Bring your own router.

## Usage
```go
import (
  "http"

  "github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

  "github.com/gorilla/mux"
  "github.com/iamatypeofwalrus/shim"
)

package main

func main() {
  r := mux.NewRouter()
  r.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
    fmt.Fprint(w, "hello, world")
  })

  shim := shim.New(r)
  lambda.Start(shim.Handle)
}
```