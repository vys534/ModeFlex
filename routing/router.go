package routing

import (
	"ModeFlex/routing/routes"
	"github.com/valyala/fasthttp"
)

func Route(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/card":
		routes.RouteCard(ctx)
		return
	default:
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	}
}
