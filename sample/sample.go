package main

import (
	"fmt"
	"github.com/amirimatin/gapBotApi"
)

func main() {
	api, err := gapBotApi.NewBotAPI("7e1e2fe810ef85e8489d5772fa297028879be1e674b4656e63951ffa7c0759da")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	api.HandleMessage("/Hi", func(botApi *gapBotApi.BotAPI, message *gapBotApi.Message) error {
		msg := gapBotApi.NewMessage(418705986, "sample")
		api.Send(msg)
		return nil
	})
	msg := gapBotApi.NewMessage(418705986, "sample messageHandler")
	msg.ReplyKeyboardMarkup = gapBotApi.NewReplyKeyboardMarkup(
		gapBotApi.NewKeyboardButtonRow(
			gapBotApi.NewKeyboardButton("YES", "yes"),
			gapBotApi.NewKeyboardButton("NO", "no"),
		),
		gapBotApi.NewKeyboardButtonRow(gapBotApi.NewKeyboardButton("CANCEL", "cancel")),
		gapBotApi.NewKeyboardButtonRow(gapBotApi.NewKeyboardButtonLocation("مکان شما")),
		gapBotApi.NewKeyboardButtonRow(gapBotApi.NewKeyboardButtonContact("تلفن شما")),
	)
	msg.InlineKeyboardMarkup = gapBotApi.NewInlineKeyboardMarkup(gapBotApi.NewInlineKeyboardRow(gapBotApi.NewInlineKeyboardButton("hi", gapBotApi.CallbackQueryAction{
		StatePath: "salam",
		Params: map[string]string{
			"name": "amiri",
		},
	}), gapBotApi.NewInlineKeyboardButtonURL("google", "https://google.com", gapBotApi.INLINE_KEYBOARD_URL_OPENIN_WEBVIEW)),
		gapBotApi.NewInlineKeyboardRow(gapBotApi.NewInlineKeyboardButtonURL("google.com", "https://google.com", gapBotApi.INLINE_KEYBOARD_URL_OPENIN_WEBVIEW_FULL)),
		gapBotApi.NewInlineKeyboardRow(gapBotApi.NewInlineKeyboardButtonURL("google.com", "https://google.com", gapBotApi.INLINE_KEYBOARD_URL_OPENIN_WEBVIEW_WITH_HEADER)),
		gapBotApi.NewInlineKeyboardRow(gapBotApi.NewInlineKeyboardButtonPayment("پرداخت کنید", 100, gapBotApi.INLINE_KEYBOARD_CURRENCY_IRR, "11454654sfdf5gv4d56144212", "همبنجوری الکی")))

	options := []gapBotApi.FormObjectOption{
		{"male": "male"},
		{"fmale": "fmale"},
	}
	msg.Form = gapBotApi.NewForm(gapBotApi.NewFormObjectQrcode("scan", "scan me"), gapBotApi.NewFormObjectCheckbox("agrre", "I Agree"),
		gapBotApi.NewFormObjectRadioInput("gender", "gender", options),
		gapBotApi.NewFormObjectSelect("gender", "gender", options),
		gapBotApi.NewFormObjectSubmit("ارسال", "ارسال"),
	)
	//photo := gapBotApi.FilePath("/home/amiri/Pictures/background/216_20160615_1074822377.jpg")
	//mPhoto := gapBotApi.NewPhoto(415661068, photo)
	//mPhoto.Description = "my image background"
	//fmt.Println(api.Send(mPhoto))

	//video := gapBotApi.FilePath("/home/amiri/Downloads/sample.json")
	//mVideo := gapBotApi.NewFile(415661068, video)
	//mVideo.Description = "sample json file"
	//fmt.Println(api.Send(mVideo))
	fmt.Println(api.Send(msg))
}
