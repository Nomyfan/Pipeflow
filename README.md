# Pipeflow
Pipeflow is a middleware container which is planned to used in my own blog system.

## Quick Start
```golang
import (
	"fmt"
	"net/http"
	"pipeflow/core"
	Pipeflow "pipeflow/flow"
)

func main() {
	flow := Pipeflow.NewFlow()
	flow.UseCors([]string{"http://localhost:18080"}, nil, nil, nil)
	_ = flow.Register("/api/index/greet", apiGreet, []core.HttpMethod{core.HttpGet})
	_ = flow.Register("/hello/{foo}/{bar}/tail?uid=?&name=?", helloHandler, []core.HttpMethod{core.HttpGet, core.HttpPost})
	flow.Run(loggerMiddleware{})
	flow.Use(tokenChecker{})
	_ = http.ListenAndServe(":12080", flow)
}

func apiGreet(ctx core.HttpContext) {
	fmt.Println(ctx.Request.Host)
	_, _ = fmt.Fprintln(ctx.ResponseWriter, "1")
}

func helloHandler(ctx core.HttpContext) {
	_, _ = fmt.Fprintln(ctx.ResponseWriter, "<h1>Pipeflow</h1></br> Foo: "+(*ctx.Vars)["foo"]+"</br> Bar: "+(*ctx.Vars)["bar"])
}

type loggerMiddleware struct {
}

func (lmw loggerMiddleware) Handle(ctx core.HttpContext) {
	fmt.Println("Request from " + ctx.Request.Header.Get("Origin"))
	fmt.Println("The path is " + ctx.Request.URL.Path)
	fmt.Println("Http method is " + ctx.Request.Method)
}

type tokenChecker struct {
}

func (tc tokenChecker) Handle(ctx core.HttpContext) bool {
	if "" != ctx.Request.Header.Get("token") {
		return true
	}

	fmt.Println("Cannot access")
	return false
}
```