package http_client

import (
	"ModeFlex/api"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"time"
)

func JSONResponse(ctx *fasthttp.RequestCtx, statusCode int, body *api.ServerResponse) {
	if statusCode == 1 {
		fmt.Printf("[ERROR - %s] Status code %d, message: %s\n", time.Now().UTC().Format(time.RFC822), statusCode, body.Message)
	}
	ctx.Response.Header.SetCanonical([]byte("Content-Type"), []byte("application/json"))
	ctx.Response.SetStatusCode(statusCode)

	if err := json.NewEncoder(ctx).Encode(body); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}
