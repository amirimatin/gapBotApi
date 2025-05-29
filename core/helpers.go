// helpers.go
package gapBotApi

import (
	"encoding/json"

	"github.com/amirimatin/gapBotApi/v2/models"
)

func NewMessage(chatID int64, text string) MessageConfig {
	return MessageConfig{
		BaseChat: BaseChat{ChatID: chatID},
		Text:     text,
		Type:     MESSAGE_TYPE_TEXT,
	}
}

func NewUpdateMessage(chatID int64, messageID int64, text string) UpdateMessageConfig {
	return UpdateMessageConfig{
		BaseChat:  BaseChat{ChatID: chatID},
		MessageId: messageID,
		Text:      text,
		Type:      MESSAGE_TYPE_TEXT,
	}
}

func NewDeleteMessage(chatID int64, messageID int64) DeleteMessageConfig {
	return DeleteMessageConfig{
		BaseChat:  BaseChat{ChatID: chatID},
		MessageId: messageID,
	}
}

func NewKeyboardButton(text string, value string) models.ReplyKeyboardButton {
	return models.ReplyKeyboardButton{"text": text, "value": value}
}

func NewKeyboardButtonRow(buttons ...models.ReplyKeyboardButton) []models.ReplyKeyboardButton {
	return append([]models.ReplyKeyboardButton{}, buttons...)
}

func NewReplyKeyboardMarkup(rows ...[]models.ReplyKeyboardButton) models.ReplyKeyboardMarkup {
	return models.ReplyKeyboardMarkup{Keyboard: append([][]models.ReplyKeyboardButton{}, rows...)}
}

func NewKeyboardButtonLocation(text string) models.ReplyKeyboardButton {
	return models.ReplyKeyboardButton{"$location": text}
}

func NewKeyboardButtonContact(text string) models.ReplyKeyboardButton {
	return models.ReplyKeyboardButton{"$contact": text}
}

func NewInlineKeyboardRow(buttons ...models.InlineKeyboardButton) []models.InlineKeyboardButton {
	return append([]models.InlineKeyboardButton{}, buttons...)
}

func NewInlineKeyboardMarkup(rows ...[]models.InlineKeyboardButton) models.InlineKeyboardMarkup {
	return append([][]models.InlineKeyboardButton{}, rows...)
}

func NewInlineKeyboardButton(text string, action models.CallbackQueryAction) models.InlineKeyboardButton {
	btn := models.InlineKeyboardButton{Text: text}
	if jsonData, err := json.Marshal(action); err == nil {
		btn.CallbackData = string(jsonData)
	}
	return btn
}

func NewInlineKeyboardButtonURL(text, url, openIn string) models.InlineKeyboardButton {
	return models.InlineKeyboardButton{
		Text:   text,
		URL:    url,
		OpenIn: openIn,
	}
}

func NewInlineKeyboardButtonPayment(text string, amount int, currency, refId, desc string) models.InlineKeyboardButton {
	return models.InlineKeyboardButton{
		Text:        text,
		Amount:      amount,
		Currency:    currency,
		RefId:       refId,
		Description: desc,
	}
}

func NewFormObjectInput(name, label string, value ...string) models.FormObject {
	val := ""
	if len(value) > 0 {
		val = value[0]
	}
	return models.FormObject{Name: name, Label: label, Value: val, Type: FORM_OBJECTS_TYPE_TEXT}
}

func NewFormObjectTextarea(name, label string, value ...string) models.FormObject {
	val := ""
	if len(value) > 0 {
		val = value[0]
	}
	return models.FormObject{Name: name, Label: label, Value: val, Type: FORM_OBJECTS_TYPE_TEXTAREA}
}

func NewFormObjectCheckbox(name, label string) models.FormObject {
	return models.FormObject{Name: name, Label: label, Type: FORM_OBJECTS_TYPE_CHECKBOX}
}

func NewFormObjectBarcode(name, label string) models.FormObject {
	return models.FormObject{Name: name, Label: label, Value: "barcode", Type: FORM_OBJECTS_TYPE_INBUILT}
}

func NewFormObjectQrcode(name, label string) models.FormObject {
	return models.FormObject{Name: name, Label: label, Value: "qrcode", Type: FORM_OBJECTS_TYPE_INBUILT}
}

func NewFormObjectSubmit(name, label string) models.FormObject {
	return models.FormObject{Name: name, Label: label, Type: FORM_OBJECTS_TYPE_SUBMIT}
}

func NewFormObjectInputWithValue(name, label, value string) models.FormObject {
	return models.FormObject{Name: name, Label: label, Value: value, Type: FORM_OBJECTS_TYPE_TEXT}
}

func NewFormObjectRadioInput(name, label string, options []models.FormObjectOption) models.FormObject {
	return models.FormObject{Name: name, Label: label, Type: FORM_OBJECTS_TYPE_RADIO, Options: options}
}

func NewFormObjectSelect(name, label string, options []models.FormObjectOption) models.FormObject {
	return models.FormObject{Name: name, Label: label, Type: FORM_OBJECTS_TYPE_SELECT, Options: options}
}

func NewForm(objects ...models.FormObject) []models.FormObject {
	return append([]models.FormObject{}, objects...)
}

func NewPhoto(chatID int64, file RequestFileData) PhotoConfig {
	return PhotoConfig{BaseFile: BaseFile{BaseChat: BaseChat{ChatID: chatID}, File: file}}
}

func NewVideo(chatID int64, file RequestFileData) VideoConfig {
	return VideoConfig{BaseFile: BaseFile{BaseChat: BaseChat{ChatID: chatID}, File: file}}
}

func NewVoice(chatID int64, file RequestFileData) VoiceConfig {
	return VoiceConfig{BaseFile: BaseFile{BaseChat: BaseChat{ChatID: chatID}, File: file}}
}

func NewAudio(chatID int64, file RequestFileData) AudioConfig {
	return AudioConfig{BaseFile: BaseFile{BaseChat: BaseChat{ChatID: chatID}, File: file}}
}

func NewFile(chatID int64, file RequestFileData) FileConfig {
	return FileConfig{BaseFile: BaseFile{BaseChat: BaseChat{ChatID: chatID}, File: file}}
}

func NewAnswerCallback(chatID int64, callbackId, text string, showAlert bool) CallbackAnswerConfig {
	return CallbackAnswerConfig{
		BaseChat:   BaseChat{ChatID: chatID},
		CallbackId: callbackId,
		Text:       text,
		ShowAlert:  showAlert,
	}
}
