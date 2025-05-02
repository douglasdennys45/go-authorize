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

type responseData struct {
	IsAuthorized bool                   `json:"isAuthorized"`
	HashMap      string                 `json:"hashMap"`
	Who          map[string]interface{} `json:"who"`
}

type responses struct {
	RequestId string      `json:"requestId"`
	Timestamp string      `json:"timestamp"`
	Data      interface{} `json:"data"`
}

func (md *middleware) Authorize(ctx *fiber.Ctx) error {
	req, err := http.NewRequest("POST", md.url, nil)
	if err != nil {
		return response.RenderJSON(ctx, err.Error(), 403)
	}
	if ctx.GetReqHeaders()["X-Access-Token-Api"] == nil {
		return response.RenderJSON(ctx, "Unauthorized", 403)
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
		return response.RenderJSON(ctx, "Unauthorized", 403)
	}

	responseMap, ok := data.Data.(map[string]interface{})
	if !ok {
		return response.RenderJSON(ctx, "Invalid response format", 403)
	}

	jsonData, err := json.Marshal(responseMap)
	if err != nil {
		return response.RenderJSON(ctx, "Error processing response", 500)
	}

	var responseData responseData
	if err := json.Unmarshal(jsonData, &responseData); err != nil {
		return response.RenderJSON(ctx, "Invalid response format", 403)
	}

	if !responseData.IsAuthorized {
		return response.RenderJSON(ctx, "Unauthorized", 403)
	}

	ctx.Set("X-Revision-HashMap", responseData.HashMap)
	whoJSON, err := json.Marshal(responseData.Who)
	if err != nil {
		return response.RenderJSON(ctx, "Error converting who to JSON", 500)
	}
	ctx.Request().Header.Set("X-Who", string(whoJSON))

	return ctx.Next()
}
