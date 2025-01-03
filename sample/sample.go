package main

import (
	"fmt"
	"github.com/amirimatin/gapBotApi/v2"
	"os"
)

func main() {
	api, err := gapBotApi.NewBotAPI(os.Getenv("GAPBOT_TOKEN"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	api.Handle("/test", func(ctx *gapBotApi.Ctx) (gapBotApi.Message, error) {
		msg := gapBotApi.NewMessage(ctx.Message.ChatID, "hi you can start")
		msg.InlineKeyboardMarkup = gapBotApi.NewInlineKeyboardMarkup(gapBotApi.NewInlineKeyboardRow(
			gapBotApi.NewInlineKeyboardButton("step one", gapBotApi.CallbackQueryAction{
				StatePath: "/answer",
				Params:    nil,
			}),
		))
		return api.Send(msg)
	})

	api.Handle("/answer", func(ctx *gapBotApi.Ctx) (gapBotApi.Message, error) {
		msg := gapBotApi.NewAnswerCallback(ctx.Message.ChatID, ctx.Message.CallbackQuery.CallbackId, "test with show alert = false", false)
		ctx.Bot().Send(msg)
		return ctx.Next()
	}, func(ctx *gapBotApi.Ctx) (gapBotApi.Message, error) {
		msg := gapBotApi.NewAnswerCallback(ctx.Message.ChatID, ctx.Message.CallbackQuery.CallbackId, "test with show alert = true", true)
		return ctx.Bot().Send(msg)
	})

	api.Handle("/start", func(ctx *gapBotApi.Ctx) (gapBotApi.Message, error) {
		msg := gapBotApi.NewMessage(ctx.Message.ChatID, "hi you can start")
		msg.InlineKeyboardMarkup = gapBotApi.NewInlineKeyboardMarkup(gapBotApi.NewInlineKeyboardRow(
			gapBotApi.NewInlineKeyboardButton("step one", gapBotApi.CallbackQueryAction{
				StatePath: "/step1",
				Params:    nil,
			}),
		))
		return api.Send(msg)
	})

	api.Handle("/start", func(ctx *gapBotApi.Ctx) (gapBotApi.Message, error) {
		msg := gapBotApi.NewMessage(ctx.Message.ChatID, "hi you can start")
		msg.InlineKeyboardMarkup = gapBotApi.NewInlineKeyboardMarkup(gapBotApi.NewInlineKeyboardRow(
			gapBotApi.NewInlineKeyboardButton("step one", gapBotApi.CallbackQueryAction{
				StatePath: "/step1",
				Params:    nil,
			}),
		))
		return api.Send(msg)
	})

	api.Handle("/step1", func(ctx *gapBotApi.Ctx) (gapBotApi.Message, error) {
		msg := gapBotApi.NewMessage(ctx.Message.ChatID, "You in step one")
		msg.InlineKeyboardMarkup = gapBotApi.NewInlineKeyboardMarkup(gapBotApi.NewInlineKeyboardRow(
			gapBotApi.NewInlineKeyboardButton("step two", gapBotApi.CallbackQueryAction{
				StatePath: "/step2",
				Params:    nil,
			}),
		))
		return api.Send(msg)
	})
	api.Handle("/step2", func(ctx *gapBotApi.Ctx) (gapBotApi.Message, error) {
		msg := gapBotApi.NewMessage(ctx.Message.ChatID, "You in step two")
		msg.InlineKeyboardMarkup = gapBotApi.NewInlineKeyboardMarkup(gapBotApi.NewInlineKeyboardRow(
			gapBotApi.NewInlineKeyboardButton("step three", gapBotApi.CallbackQueryAction{
				StatePath: "/step3",
				Params:    nil,
			}),
		))
		return api.Send(msg)
	})
	api.Handle("/step3", func(ctx *gapBotApi.Ctx) (gapBotApi.Message, error) {
		msg := gapBotApi.NewMessage(ctx.Message.ChatID, "You in step three")
		msg.InlineKeyboardMarkup = gapBotApi.NewInlineKeyboardMarkup(gapBotApi.NewInlineKeyboardRow(
			gapBotApi.NewInlineKeyboardButton("step fore", gapBotApi.CallbackQueryAction{
				StatePath: "/step4",
				Params:    nil,
			}),
		))
		return api.Send(msg)
	})
	api.Handle("/step4", func(ctx *gapBotApi.Ctx) (gapBotApi.Message, error) {
		msg := gapBotApi.NewMessage(ctx.Message.ChatID, "You in step fore")
		msg.InlineKeyboardMarkup = gapBotApi.NewInlineKeyboardMarkup(gapBotApi.NewInlineKeyboardRow(
			gapBotApi.NewInlineKeyboardButton("step five", gapBotApi.CallbackQueryAction{
				StatePath: "/step5",
				Params:    nil,
			}),
		))
		return api.Send(msg)
	})
	api.Handle("/step5", func(ctx *gapBotApi.Ctx) (gapBotApi.Message, error) {
		msg := gapBotApi.NewMessage(ctx.Message.ChatID, "You in step five")
		msg.InlineKeyboardMarkup = gapBotApi.NewInlineKeyboardMarkup(gapBotApi.NewInlineKeyboardRow(
			gapBotApi.NewInlineKeyboardButton("Back", gapBotApi.CallbackQueryAction{
				StatePath: "/back",
				Params:    nil,
			}),
		))
		return api.Send(msg)
	})

	//msg := gapBotApi.NewMessage(433221574, "sample messageHandler")
	//send, err := api.Send(msg)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//fmt.Println(send)
	//msg.ReplyKeyboardMarkup = gapBotApi.NewReplyKeyboardMarkup(
	//	gapBotApi.NewKeyboardButtonRow(
	//		gapBotApi.NewKeyboardButton("YES", "yes"),
	//		gapBotApi.NewKeyboardButton("NO", "no"),
	//	),
	//	gapBotApi.NewKeyboardButtonRow(gapBotApi.NewKeyboardButton("CANCEL", "cancel")),
	//	gapBotApi.NewKeyboardButtonRow(gapBotApi.NewKeyboardButtonLocation("مکان شما")),
	//	gapBotApi.NewKeyboardButtonRow(gapBotApi.NewKeyboardButtonContact("تلفن شما")),
	//)
	//msg.InlineKeyboardMarkup = gapBotApi.NewInlineKeyboardMarkup(gapBotApi.NewInlineKeyboardRow(gapBotApi.NewInlineKeyboardButton("hi", gapBotApi.CallbackQueryAction{
	//	StatePath: "salam",
	//	Params: map[string]string{
	//		"name": "amiri",
	//	},
	//}), gapBotApi.NewInlineKeyboardButtonURL("google", "https://google.com", gapBotApi.INLINE_KEYBOARD_URL_OPENIN_WEBVIEW)),
	//	gapBotApi.NewInlineKeyboardRow(gapBotApi.NewInlineKeyboardButtonURL("google.com", "https://google.com", gapBotApi.INLINE_KEYBOARD_URL_OPENIN_WEBVIEW_FULL)),
	//	gapBotApi.NewInlineKeyboardRow(gapBotApi.NewInlineKeyboardButtonURL("google.com", "https://google.com", gapBotApi.INLINE_KEYBOARD_URL_OPENIN_WEBVIEW_WITH_HEADER)),
	//	gapBotApi.NewInlineKeyboardRow(gapBotApi.NewInlineKeyboardButtonPayment("پرداخت کنید", 100, gapBotApi.INLINE_KEYBOARD_CURRENCY_IRR, "11454654sfdf5gv4d56144212", "همبنجوری الکی")))

	//options := []gapBotApi.FormObjectOption{
	//	{"male": "male"},
	//	{"fmale": "fmale"},
	//}
	////
	//photo := gapBotApi.FilePath("/home/amiri/Pictures/2bd57f01e8c8f40e05656954a2271799.jpg")
	//mPhoto := gapBotApi.NewFile(433221574, photo)
	//mPhoto.Description = "my image background"
	//fmt.Println(api.Send(mPhoto))
	////
	//video := gapBotApi.FilePath("/home/amiri/Pictures/5825830340513503652.mp4")
	//mVideo := gapBotApi.NewVideo(433221574, video)
	//mVideo.Description = "sample video file"
	//fmt.Println(api.Send(mVideo))
	//fmt.Println(api.Send(msg))

	//fmt.Println(api.Send(msg))
	//fmt.Println(api.Send(mVideo))
	api.Serve(3900, "/bot/callback")

}
