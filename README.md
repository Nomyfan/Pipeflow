# Pipeflow
Pipeflow is a middleware container which is used in my own blog system.

## Quick Look
```golang
func main() {
	flow := pipeflow.NewFlow()

	flow.Run(func(ctx pipeflow.HTTPContext) {
		fmt.Println("request URL: " + ctx.Request.RequestURI)
	})

	flow.UseCors([]string{"http://localhost:18080"}, nil, nil, nil)

	flow.Use(func(ctx pipeflow.HTTPContext, next func()) {
		fmt.Println("first")
		next()
		fmt.Println("first post action")
	})

	flow.Use(func(ctx pipeflow.HTTPContext, next func()) {
		fmt.Println("second")
		next()
		fmt.Println("second post action")
	})

	flow.Use(func(ctx pipeflow.HTTPContext, next func()) {
		if token := ctx.Request.Header.Get("token"); token != "" {
			next()
		} else {
			ctx.ResponseWriter.WriteHeader(http.StatusNonAuthoritativeInfo)
			_, _ = ctx.ResponseWriter.Write([]byte("NonAuthoritativeInfo"))
		}
	})

	redisClient := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379", Password: "password", DB: 0})
	if _, err := redisClient.Ping().Result(); nil != err {
		panic(err)
	}

	flow.SetResource("redis", redisClient)

	_ = flow.Map("/hey", func(ctx pipeflow.HTTPContext) {
		var client, _ = ctx.GetResource("redis").(*redis.Client)
		var count, _ = client.Get("count").Int()
		client.Set("count", count+1, -1)
		_, _ = ctx.ResponseWriter.Write([]byte("hello"))
	}, []pipeflow.HTTPMethod{pipeflow.HTTPPost, pipeflow.HTTPGet})

	_ = flow.GET("/hello", func(ctx pipeflow.HTTPContext) {
		var client, _ = ctx.GetResource("redis").(*redis.Client)
		var count, _ = client.Get("count").Int()
		client.Set("count", count+1, -1)
		_, _ = ctx.ResponseWriter.Write([]byte("hello"))
	})

	_ = flow.POST("/{foo}/hello?id=?&name=?", func(ctx pipeflow.HTTPContext) {
		_, _ = fmt.Fprintln(ctx.ResponseWriter, "foo = "+ctx.Vars["foo"]+", id = "+ctx.Request.Form.Get("id")+", name = "+ctx.Request.Form.Get("name"))
	})

	_ = http.ListenAndServe(":8888", flow)
}
```

Request: `http://localhost:8888/bar/hello?id=1&name=nomyfan`

Response: `foo = bar, id = 1, name = nomyfan`

Console output:
```
request URL: /bar/hello?id=1&name=nomyfan
first
second
second post action
first post action
```