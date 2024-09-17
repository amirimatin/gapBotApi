## gapBotApi Go Package
This package is a Go client library for interacting with GAP Messenger's Bot API. It simplifies sending and receiving messages, as well as handling various types of keyboard inputs, media, and forms.

### Installation
To install this package, use go get:
``` go
go get github.com/amirimatin/gapBotApi
```
### Usage
Here's a sample usage of the package:

``` go
package main

import (
	"fmt"
	"github.com/amirimatin/gapBotApi"
)

func main() {
	// Create a new bot API instance with your token
	api, err := gapBotApi.NewBotAPI("your_bot_api_token_here")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Handle incoming messages with the "/Hi" command
	api.HandleMessage("/Hi", func(botApi *gapBotApi.BotAPI, message *gapBotApi.Message) error {
		msg := gapBotApi.NewMessage(418705986, "sample response")
		api.Send(msg)
		return nil
	})
    api.HandleCallback("hello", AcceptJoinRequest)
    
    
	// Create a new message with keyboard options
	msg := gapBotApi.NewMessage(418705986, "sample messageHandler")

	// Define a reply keyboard with buttons
	msg.ReplyKeyboardMarkup = gapBotApi.NewReplyKeyboardMarkup(
		gapBotApi.NewKeyboardButtonRow(
			gapBotApi.NewKeyboardButton("YES", "yes"),
			gapBotApi.NewKeyboardButton("NO", "no"),
		),
		gapBotApi.NewKeyboardButtonRow(gapBotApi.NewKeyboardButton("CANCEL", "cancel")),
		gapBotApi.NewKeyboardButtonRow(gapBotApi.NewKeyboardButtonLocation("Your Location")),
		gapBotApi.NewKeyboardButtonRow(gapBotApi.NewKeyboardButtonContact("Your Contact")),
	)

	// Define an inline keyboard with URLs and actions
	msg.InlineKeyboardMarkup = gapBotApi.NewInlineKeyboardMarkup(
		gapBotApi.NewInlineKeyboardRow(
			gapBotApi.NewInlineKeyboardButton("hi", gapBotApi.CallbackQueryAction{
				StatePath: "hello",
				Params: map[string]string{
					"name": "amiri",
				},
			}),
			gapBotApi.NewInlineKeyboardButtonURL("Google", "https://google.com", gapBotApi.INLINE_KEYBOARD_URL_OPENIN_WEBVIEW),
		),
		gapBotApi.NewInlineKeyboardRow(gapBotApi.NewInlineKeyboardButtonPayment("Make Payment", 100, gapBotApi.INLINE_KEYBOARD_CURRENCY_IRR, "payment_id", "Payment Description")),
	)

	// Define a form with input options
	options := []gapBotApi.FormObjectOption{
		{"male": "male"},
		{"female": "female"},
	}
	msg.Form = gapBotApi.NewForm(
		gapBotApi.NewFormObjectQrcode("scan", "Scan this code"),
		gapBotApi.NewFormObjectCheckbox("agree", "I Agree"),
		gapBotApi.NewFormObjectRadioInput("gender", "Gender", options),
		gapBotApi.NewFormObjectSelect("gender", "Gender", options),
		gapBotApi.NewFormObjectSubmit("Submit", "Submit"),
	)

	// Send a photo
	photo := gapBotApi.FilePath("/path/to/photo.jpg")
	mPhoto := gapBotApi.NewPhoto(418705986, photo)
	mPhoto.Description = "Sample Photo"
	fmt.Println(api.Send(mPhoto))

	// Send a video or any other file
	video := gapBotApi.FilePath("/path/to/file.mp4")
	mVideo := gapBotApi.NewFile(418705986, video)
	mVideo.Description = "Sample Video"
	fmt.Println(api.Send(mVideo))

	// Send the message with keyboards and form
	fmt.Println(api.Send(msg))
}
func AcceptJoinRequest(botApi *gapBotApi.BotAPI, callback *gapBotApi.CallbackQuery) error {
		// Process accepting join request
		msg := gapBotApi.NewMessage(callback.From.ID, "Join request accepted")
		return botApi.Send(msg)
}
```

### Features
* Handle Messages: Easily handle messages based on commands or text.
* Handle Callbacks: Capture and respond to user interactions like button clicks using callback handlers.
* Reply and Inline Keyboards: Create and send interactive keyboards.
* Forms: Use forms with inputs like checkboxes, radio buttons, and file uploads.
* Media Support: Send photos, videos, and files.
* Payment Integration: Add inline buttons for making payments.

### Example Callback Handling

``` go
api.HandleCallback("admin.join.accept", AcceptJoinRequest)

func AcceptJoinRequest(botApi *gapBotApi.BotAPI, callback *gapBotApi.CallbackQuery) error {
	msg := gapBotApi.NewMessage(callback.From.ID, "Join request accepted")
	return botApi.Send(msg)
}
```
### License

This project is licensed under the MIT License. This means you are free to use, modify, distribute, and incorporate this project in your own software, even for commercial purposes, as long as you include the original copyright notice and this license in any substantial portions of the software.

The software is provided "as is", without warranty of any kind. For more details, refer to the LICENSE file.