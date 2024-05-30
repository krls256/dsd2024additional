package http

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func QueryToMap(ctx *fiber.Ctx) map[string]interface{} {
	paramMap := make(map[string]interface{}, 0)

	for k, v := range ctx.Queries() {
		paramMap[k] = v

		i, err := strconv.Atoi(v)
		if err == nil {
			paramMap[k] = i
		}
	}

	return paramMap
}

func ParseRequest[T any](payload interface{}, req T) error {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, &req); err != nil {
		return err
	}

	return nil
}
