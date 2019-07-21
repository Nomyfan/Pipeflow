# Pipeflow
Pipeflow is middleware container which is planned to used in my own blog system.

## Example
```golang
import (
	"fmt"
	"net/http"
	"pipeflow/core"
	Pipeflow "pipeflow/flow"
)

func main() {
	flow := Pipeflow.NewFlow()
	flow.AddHandler("/hello", helloHandler, []core.HttpMethod{core.HttpGet, core.HttpPost})
	flow.Run(LoggerMiddleware{})
	flow.Use(TokenChecker{})
	_ = http.ListenAndServe(":8080", flow)
}

func helloHandler(ctx core.HttpContext) {
	_, _ = fmt.Fprintln(ctx.ResponseWriter, "<h1>Pipeflow</h1>")
}

type LoggerMiddleware struct {
}

func (lmw LoggerMiddleware) Handle(ctx core.HttpContext) {
	fmt.Println("Request from " + ctx.Request.Host)
	fmt.Println("The path is " + ctx.Request.URL.Path)
	fmt.Println("Http method is " + ctx.Request.Method)
}

type TokenChecker struct {
}

func (tc TokenChecker) Handle(ctx core.HttpContext) bool {
	if "" != ctx.Request.Header.Get("token") {
		return true
    }
    
	fmt.Println("Cannot access")
	return false
}
```

## TODO
- One handler binds several URL
- Support fetching params from URL