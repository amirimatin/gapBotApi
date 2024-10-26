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
	api.Handle("/test", func(ctx *gapBotApi.Ctx) error {
		msg := gapBotApi.NewMessage(ctx.Message.ChatID, "hi you can start")
		msg.InlineKeyboardMarkup = gapBotApi.NewInlineKeyboardMarkup(gapBotApi.NewInlineKeyboardRow(
			gapBotApi.NewInlineKeyboardButton("step one", gapBotApi.CallbackQueryAction{
				StatePath: "/answer",
				Params:    nil,
			}),
		))
		send, err := api.Send(msg)
		fmt.Println(err, send)
		return nil
	})

	api.Handle("/answer", func(ctx *gapBotApi.Ctx) error {
		msg := gapBotApi.NewAnswerCallback(ctx.Message.ChatID, ctx.Message.CallbackQuery.CallbackId, "test with show alert = false", false)
		ctx.Bot().Send(msg)
		return ctx.Next()
	}, func(ctx *gapBotApi.Ctx) error {
		msg := gapBotApi.NewAnswerCallback(ctx.Message.ChatID, ctx.Message.CallbackQuery.CallbackId, "test with show alert = true", true)
		ctx.Bot().Send(msg)
		return nil
	})

	api.Handle("/start", func(ctx *gapBotApi.Ctx) error {
		msg := gapBotApi.NewMessage(ctx.Message.ChatID, "hi you can start")
		msg.InlineKeyboardMarkup = gapBotApi.NewInlineKeyboardMarkup(gapBotApi.NewInlineKeyboardRow(
			gapBotApi.NewInlineKeyboardButton("step one", gapBotApi.CallbackQueryAction{
				StatePath: "/step1",
				Params:    nil,
			}),
		))
		send, err := api.Send(msg)
		fmt.Println(err, send)
		return nil
	})

	api.Handle("/start", func(ctx *gapBotApi.Ctx) error {
		msg := gapBotApi.NewMessage(ctx.Message.ChatID, "hi you can start")
		msg.InlineKeyboardMarkup = gapBotApi.NewInlineKeyboardMarkup(gapBotApi.NewInlineKeyboardRow(
			gapBotApi.NewInlineKeyboardButton("step one", gapBotApi.CallbackQueryAction{
				StatePath: "/step1",
				Params:    nil,
			}),
		))
		send, err := api.Send(msg)
		fmt.Println(err, send)
		return nil
	})

	api.Handle("/step1", func(ctx *gapBotApi.Ctx) error {
		msg := gapBotApi.NewMessage(ctx.Message.ChatID, "You in step one")
		msg.InlineKeyboardMarkup = gapBotApi.NewInlineKeyboardMarkup(gapBotApi.NewInlineKeyboardRow(
			gapBotApi.NewInlineKeyboardButton("step two", gapBotApi.CallbackQueryAction{
				StatePath: "/step2",
				Params:    nil,
			}),
		))
		send, err := api.Send(msg)
		fmt.Println(err, send)
		return nil
	})
	api.Handle("/step2", func(ctx *gapBotApi.Ctx) error {
		msg := gapBotApi.NewMessage(ctx.Message.ChatID, "You in step two")
		msg.InlineKeyboardMarkup = gapBotApi.NewInlineKeyboardMarkup(gapBotApi.NewInlineKeyboardRow(
			gapBotApi.NewInlineKeyboardButton("step three", gapBotApi.CallbackQueryAction{
				StatePath: "/step3",
				Params:    nil,
			}),
		))
		send, err := api.Send(msg)
		fmt.Println(err, send)
		return nil
	})
	api.Handle("/step3", func(ctx *gapBotApi.Ctx) error {
		msg := gapBotApi.NewMessage(ctx.Message.ChatID, "You in step three")
		msg.InlineKeyboardMarkup = gapBotApi.NewInlineKeyboardMarkup(gapBotApi.NewInlineKeyboardRow(
			gapBotApi.NewInlineKeyboardButton("step fore", gapBotApi.CallbackQueryAction{
				StatePath: "/step4",
				Params:    nil,
			}),
		))
		send, err := api.Send(msg)
		fmt.Println(err, send)
		return nil
	})
	api.Handle("/step4", func(ctx *gapBotApi.Ctx) error {
		msg := gapBotApi.NewMessage(ctx.Message.ChatID, "You in step fore")
		msg.InlineKeyboardMarkup = gapBotApi.NewInlineKeyboardMarkup(gapBotApi.NewInlineKeyboardRow(
			gapBotApi.NewInlineKeyboardButton("step five", gapBotApi.CallbackQueryAction{
				StatePath: "/step5",
				Params:    nil,
			}),
		))
		send, err := api.Send(msg)
		fmt.Println(err, send)
		return nil
	})
	api.Handle("/step5", func(ctx *gapBotApi.Ctx) error {
		msg := gapBotApi.NewMessage(ctx.Message.ChatID, "You in step five")
		msg.InlineKeyboardMarkup = gapBotApi.NewInlineKeyboardMarkup(gapBotApi.NewInlineKeyboardRow(
			gapBotApi.NewInlineKeyboardButton("Back", gapBotApi.CallbackQueryAction{
				StatePath: "/back",
				Params:    nil,
			}),
		))
		send, err := api.Send(msg)
		fmt.Println(err, send)
		return nil
	})
	api.Serve(3900, "/bot/callback")
	//msg := gapBotApi.NewMessage(418705986, "sample messageHandler")
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
	//
	//options := []gapBotApi.FormObjectOption{
	//	{"male": "male"},
	//	{"fmale": "fmale"},
	//}

	//photo := gapBotApi.FilePath("/home/amiri/Pictures/background/216_20160615_1074822377.jpg")
	//mPhoto := gapBotApi.NewPhoto(415661068, photo)
	//mPhoto.Description = "my image background"
	//fmt.Println(api.Send(mPhoto))
	//
	//video := gapBotApi.FilePath("/home/amiri/Downloads/sample.json")
	//mVideo := gapBotApi.NewFile(415661068, video)
	//mVideo.Description = "sample json file"
	//fmt.Println(api.Send(mVideo))
	//fmt.Println(api.Send(msg))
}
