package middleware

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/douglasdennys45/go-shared/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type middleware struct {
	url string
}

func NewAuthorizeAdapter(url string) *middleware {
	return &middleware{url}
}

type responses struct {
	RequestId string `json:"requestId"`
	Timestamp string `json:"timestamp"`
	Data      bool   `json:"data"`
}

func (md *middleware) Authorize(ctx *fiber.Ctx) error {
	req, err := http.NewRequest("POST", md.url, nil)
	if err != nil {
		return response.RenderJSON(ctx, err.Error(), 403)
	}
	if ctx.GetReqHeaders()["X-Access-Token-Api"] == nil {
		return response.RenderJSON(ctx, "Não autorizado", 403)
	}
	req.Header.Set("X-Access-Token-Api", ctx.GetReqHeaders()["X-Access-Token-Api"][0])
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return response.RenderJSON(ctx, err.Error(), 403)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response.RenderJSON(ctx, err.Error(), 403)
	}
	var data responses
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		return response.RenderJSON(ctx, err.Error(), 403)
	}
	if resp.StatusCode != 200 {
		return response.RenderJSON(ctx, err.Error(), 403)
	}
	if !data.Data {
		return response.RenderJSON(ctx, "Não autorizado", 403)
	}
	return ctx.Next()
}
