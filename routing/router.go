package routing

import (
	"ModeFlex/routing/routes"
	"github.com/valyala/fasthttp"
)

func Route(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/":
		routes.RouteIndex(ctx)
		return
	default:
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	}
}
