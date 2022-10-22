package routes

import (
	"ModeFlex/api"
	"ModeFlex/data"
	"ModeFlex/drawing"
	"ModeFlex/http_client"
	"fmt"
	"github.com/valyala/fasthttp"
	"io"
	"os"
	"strconv"
	"time"
)

func RouteCard(ctx *fasthttp.RequestCtx) {
	if len(ctx.QueryArgs().Peek("id")) == 0 {
		http_client.JSONResponse(ctx, fasthttp.StatusBadRequest, &api.ServerResponse{
			Message: "No ID provided.",
		})
		return
	}

	idInt, err := strconv.ParseInt(string(ctx.QueryArgs().Peek("id")), 10, 64)
	if err != nil {
		http_client.JSONResponse(ctx, fasthttp.StatusBadRequest, &api.ServerResponse{
			Message: "Not a valid integer.",
		})
		return
	}

	var isWhitelisted bool
	for _, id := range data.Configuration.Whitelist {
		if id == idInt {
			isWhitelisted = true
		}
	}

	if !isWhitelisted {
		http_client.JSONResponse(ctx, fasthttp.StatusBadRequest, &api.ServerResponse{
			Message: "User ID not whitelisted.",
		})
		return
	}

	// Send cached image
	if time.Now().Unix()-data.LastUpdatedMap.Map[idInt] < data.Configuration.UpdateInterval {
		SendUserCard(ctx, idInt)
		return
	}

	// Else update the card
	data.LastUpdatedMap.Set(idInt, time.Now().Unix())
	err = http_client.SaveUserAvatar(idInt)
	if err != nil {
		http_client.JSONResponse(ctx, fasthttp.StatusInternalServerError, &api.ServerResponse{
			Status:  1,
			Message: fmt.Sprintf("Failed to save user avatar: %v", err),
		})
		return
	}

	userModeData, err := http_client.GetUserAllModeData(idInt)
	if err != nil {
		http_client.JSONResponse(ctx, fasthttp.StatusInternalServerError, &api.ServerResponse{
			Status:  1,
			Message: fmt.Sprintf("Failed to fetch all mode data for user ID %d: %v", idInt, err),
		})
	}

	err = drawing.CreateCard(idInt, string(ctx.QueryArgs().Peek("banner")), userModeData)
	if err != nil {
		http_client.JSONResponse(ctx, fasthttp.StatusInternalServerError, &api.ServerResponse{
			Status:  1,
			Message: fmt.Sprintf("Failed to draw card for user ID %d: %v", idInt, err),
		})
		return
	}

	SendUserCard(ctx, idInt)
}

func SendUserCard(ctx *fasthttp.RequestCtx, id int64) {

	f, err := os.Open(fmt.Sprintf("./assets/generated/user_cards/%d.png", id))
	if err != nil {
		if os.IsNotExist(err) {
			http_client.JSONResponse(ctx, fasthttp.StatusNotFound, &api.ServerResponse{
				Message: fmt.Sprintf("Card for user ID %d was not found", id),
			})
			return
		}
	}

	ctx.Response.Header.Add(fasthttp.HeaderContentType, "image/png")
	ctx.Response.Header.Add(fasthttp.HeaderCacheControl, "no-cache, no-store, must-revalidate")
	ctx.Response.Header.Add(fasthttp.HeaderPragma, "no-cache")
	ctx.Response.Header.Add(fasthttp.HeaderExpires, "0")

	_, err = io.Copy(ctx.Response.BodyWriter(), f)
	if err != nil {
		http_client.JSONResponse(ctx, fasthttp.StatusInternalServerError, &api.ServerResponse{
			Status:  1,
			Message: fmt.Sprintf(fmt.Sprintf("Failed to write file to response: %v", err)),
		})
	}
}
