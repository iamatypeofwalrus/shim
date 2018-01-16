# Lambda Routing Shim

Lambda Routing Shim provides a thin layer between the Lambda API Gateway integrations and the standard library `http.Handler` interface. Bring your own router.

```go
import (
  "http"

  "github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/gorilla/mux"
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