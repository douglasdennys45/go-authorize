package response

import (
	"time"

	"github.com/douglasdennys45/go-authorize/pkg/uuid"

	"github.com/gofiber/fiber/v2"
)

type errorDetail struct {
	Detail string `json:"detail"`
	Code   int    `json:"code"`
}

type handlerJSON struct {
	RequestId string       `json:"requestId"`
	Timestamp time.Time    `json:"timestamp"`
	Data      *interface{} `json:"data,omitempty"`
	Error     *errorDetail `json:"error,omitempty"`
	Meta      *interface{} `json:"meta,omitempty"`
}

func RenderJSON(ctx *fiber.Ctx, data interface{}, code int) error {
	if code == 403 {
		handler := handlerJSON{
			RequestId: uuid.NewUUID(),
			Timestamp: time.Now(),
			Data:      nil,
			Error:     &errorDetail{Detail: data.(string), Code: code},
		}
		return ctx.Status(code).JSON(handler)
	}
	if code == 422 {
		handler := handlerJSON{
			RequestId: uuid.NewUUID(),
			Timestamp: time.Now(),
			Data:      nil,
			Error:     &errorDetail{Detail: data.(string), Code: code},
		}
		return ctx.Status(code).JSON(handler)
	}
	handler := handlerJSON{
		RequestId: uuid.NewUUID(),
		Timestamp: time.Now(),
		Data:      &data,
		Error:     nil,
	}
	return ctx.Status(code).JSON(handler)
}

func RenderMeta(ctx *fiber.Ctx, result, meta interface{}, code int) error {
	handler := handlerJSON{
		RequestId: uuid.NewUUID(),
		Timestamp: time.Now(),
		Data:      &result,
		Error:     nil,
		Meta:      &meta,
	}
	return ctx.Status(code).JSON(handler)
}
