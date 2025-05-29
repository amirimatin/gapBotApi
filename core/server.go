package gapBotApi

import (
	"context"
	"errors"
	"fmt"

	"github.com/amirimatin/gapBotApi/v2/models"
	"github.com/gofiber/fiber/v2"
)

// HandleUpdates unmarshals the incoming update and starts processing it.
func (bot *BotAPI) HandleUpdates(update []byte) (models.Message, error) {
	ctx := Ctx{
		bot:          bot,
		Message:      &models.Message{},
		Context:      context.Background(),
		HandlerIndex: 0,
		Params:       make(map[string]interface{}),
	}

	if err := ctx.Unmarshal(update); err != nil {
		return models.Message{}, err
	}

	return ctx.Next()
}

// Serve launches a Fiber HTTP server and routes incoming callbacks.
func (bot *BotAPI) Serve(port int, callbackEndpoint string) {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"status": "error",
				"error":  err.Error(),
			})
		},
		AppName: "Gap Bot",
	})

	bot.Handle("/back", func(ctx *Ctx) (models.Message, error) {
		fmt.Printf("🌀 Back called. Stack size: %d\n", len(ctx.UserState.Stack))

		return ctx.Back()
	})

	app.Post(callbackEndpoint, func(c *fiber.Ctx) error {
		_, err := bot.HandleUpdates(c.Body())
		if err != nil {
			fmt.Printf("❌ Error handling update: %s\n", err.Error())
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"data":   nil,
		})
	})

	if err := app.Listen(fmt.Sprintf(":%d", port)); err != nil {
		panic(fmt.Sprintf("❌ Failed to start server: %v", err))
	}
}
