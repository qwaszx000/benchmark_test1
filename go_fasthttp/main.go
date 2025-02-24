package main

import (
	"github.com/valyala/fasthttp"
)

func handle_test(ctx *fasthttp.RequestCtx) {
	if string(ctx.Path()) != "/test_plain" || !ctx.IsGet() {
		ctx.SetStatusCode(404)
		return
	}

	ctx.SetBodyString("Hello world!")
}

func main() {
	fasthttp.ListenAndServe("127.0.0.1:8080", handle_test)
}
