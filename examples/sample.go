// sample.go
package main

import (
	"fmt"
	gapBotApi "github.com/amirimatin/gapBotApi/v2/core"
	"os"

	"github.com/amirimatin/gapBotApi/v2/models"
)

func main() {
	api, err := gapBotApi.NewBotAPI(os.Getenv("GAPBOT_TOKEN"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Text message with buttons
	api.Handle("/start", func(ctx *gapBotApi.Ctx) (models.Message, error) {
		msg := gapBotApi.NewMessage(ctx.Message.ChatID, "Welcome! Choose an option:")
		msg.InlineKeyboardMarkup = gapBotApi.NewInlineKeyboardMarkup(
			gapBotApi.NewInlineKeyboardRow(
				gapBotApi.NewInlineKeyboardButton("Image", models.CallbackQueryAction{StatePath: "/image"}),
				gapBotApi.NewInlineKeyboardButton("Video", models.CallbackQueryAction{StatePath: "/video"}),
			),
			gapBotApi.NewInlineKeyboardRow(
				gapBotApi.NewInlineKeyboardButton("Audio", models.CallbackQueryAction{StatePath: "/audio"}),
				gapBotApi.NewInlineKeyboardButton("Voice", models.CallbackQueryAction{StatePath: "/voice"}),
			),
			gapBotApi.NewInlineKeyboardRow(
				gapBotApi.NewInlineKeyboardButton("File", models.CallbackQueryAction{StatePath: "/file"}),
			),
		)
		return api.Send(msg)
	})

	// Send Image
	api.Handle("/image", func(ctx *gapBotApi.Ctx) (models.Message, error) {
		file := gapBotApi.FilePath("./assets/sample.jpg")
		msg := gapBotApi.NewPhoto(ctx.Message.ChatID, file)
		msg.Description = "Sample image"
		return api.Send(msg)
	})

	// Send Video
	api.Handle("/video", func(ctx *gapBotApi.Ctx) (models.Message, error) {
		file := gapBotApi.FilePath("./assets/sample.mp4")
		msg := gapBotApi.NewVideo(ctx.Message.ChatID, file)
		msg.Description = "Sample video"
		return api.Send(msg)
	})

	// Send Audio
	api.Handle("/audio", func(ctx *gapBotApi.Ctx) (models.Message, error) {
		file := gapBotApi.FilePath("./assets/sample.mp3")
		msg := gapBotApi.NewAudio(ctx.Message.ChatID, file)
		msg.Description = "Sample audio"
		return api.Send(msg)
	})

	// Send Voice
	api.Handle("/voice", func(ctx *gapBotApi.Ctx) (models.Message, error) {
		file := gapBotApi.FilePath("/home/amiri/Downloads/file_example_OOG_1MG.ogg")
		msg := gapBotApi.NewVoice(ctx.Message.ChatID, file)
		msg.Description = "Sample voice message"
		return api.Send(msg)
	})

	// Send File
	api.Handle("/file", func(ctx *gapBotApi.Ctx) (models.Message, error) {
		file := gapBotApi.FilePath("./assets/sample.pdf")
		msg := gapBotApi.NewFile(ctx.Message.ChatID, file)
		msg.Description = "Sample PDF file"
		return api.Send(msg)
	})

	api.Serve(8080, "/bot/callback")
}
