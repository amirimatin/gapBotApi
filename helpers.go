package gapBotApi

import (
	"encoding/json"
	"strconv"
)

func NewMessage(chatID int64, text string) MessageConfig {
	return MessageConfig{
		BaseChat: BaseChat{
			ChatID: chatID,
		},
		Text: text,
		Type: MESSAGE_TYPE_TEXT,
	}
}

func NewUpdateMessage(chatID int64, messageID int64, text string) UpdateMessageConfig {
	return UpdateMessageConfig{
		BaseChat: BaseChat{
			ChatID: chatID,
		},
		MessageId: messageID,
		Text:      text,
		Type:      MESSAGE_TYPE_TEXT,
	}
}

// NewDeleteMessage creates a request to delete a messageHandler.
func NewDeleteMessage(chatID int64, messageID int64) DeleteMessageConfig {
	return DeleteMessageConfig{
		BaseChat: BaseChat{
			ChatID: chatID,
		},
		MessageId: messageID,
	}
}

// NewKeyboardButton creates a regular keyboard button.
func NewKeyboardButton(text string, value string) ReplyKeyboardButton {
	return ReplyKeyboardButton{
		value: text,
	}
}
func NewKeyboardButtonRow(buttons ...ReplyKeyboardButton) []ReplyKeyboardButton {
	var row []ReplyKeyboardButton

	row = append(row, buttons...)

	return row
}

// NewReplyKeyboardMarkup creates a new regular keyboard with sane defaults.
func NewReplyKeyboardMarkup(rows ...[]ReplyKeyboardButton) ReplyKeyboardMarkup {
	var keyboard [][]ReplyKeyboardButton

	keyboard = append(keyboard, rows...)

	return ReplyKeyboardMarkup{
		Keyboard: keyboard,
	}
}

func NewKeyboardButtonLocation(text string) ReplyKeyboardButton {
	return ReplyKeyboardButton{
		"$location": text,
	}
}
func NewKeyboardButtonContact(text string) ReplyKeyboardButton {
	return ReplyKeyboardButton{
		"$contact": text,
	}
}

// NewInlineKeyboardRow creates an inline keyboard row with buttons.
func NewInlineKeyboardRow(buttons ...InlineKeyboardButton) []InlineKeyboardButton {
	var row []InlineKeyboardButton

	row = append(row, buttons...)

	return row
}

// NewInlineKeyboardMarkup creates a new inline keyboard.
func NewInlineKeyboardMarkup(rows ...[]InlineKeyboardButton) InlineKeyboardMarkup {
	var keyboard [][]InlineKeyboardButton
	return append(keyboard, rows...)
}
func (ikm InlineKeyboardMarkup) AddRow(rows []InlineKeyboardButton) InlineKeyboardMarkup {
	return append(ikm, rows)
}
func (ikm InlineKeyboardMarkup) AddButtonToEndRow(ikb InlineKeyboardButton) InlineKeyboardMarkup {
	if len(ikm) == 0 {
		ikm = ikm.AddRow(NewInlineKeyboardRow())
	}
	lastRowIndex := len(ikm) - 1
	lastRow := ikm[lastRowIndex]
	lastRow = append(lastRow, ikb)
	ikm[lastRowIndex] = lastRow
	return ikm
}
func NewInlineKeyboardButton(text string, callbackData CallbackQueryAction) InlineKeyboardButton {
	var ikb = InlineKeyboardButton{
		Text:         text,
		CallbackData: "",
	}
	jsonData, err := json.Marshal(callbackData)
	if err == nil {
		ikb.CallbackData = string(jsonData)
	}
	return ikb
}

func NewInlineKeyboardButtonURL(text, url string, openIn INLINE_KEYBOARD_URL_OPENIN) InlineKeyboardButton {
	return InlineKeyboardButton{
		Text:   text,
		URL:    url,
		OpenIn: openIn,
	}
}

func NewInlineKeyboardButtonPayment(text string, amount int, currency INLINE_KEYBOARD_CURRENCY, refId, description string) InlineKeyboardButton {
	return InlineKeyboardButton{
		Text:        text,
		Amount:      amount,
		Currency:    currency,
		RefId:       refId,
		Description: description,
	}
}

func NewFormObjectInput(name, label string, value ...string) FormObject {
	val := ""
	if len(value) != 0 {
		val = value[0]
	}
	return FormObject{
		Name:  name,
		Label: label,
		Value: val,
		Type:  FORM_OBJECTS_TYPE_TEXT,
	}
}
func NewFormObjectTextarea(name, label string, value ...string) FormObject {
	val := ""
	if len(value) != 0 {
		val = value[0]
	}
	return FormObject{
		Name:  name,
		Label: label,
		Value: val,
		Type:  FORM_OBJECTS_TYPE_TEXTAREA,
	}
}
func NewFormObjectCheckbox(name, label string) FormObject {
	return FormObject{
		Name:  name,
		Label: label,
		Type:  FORM_OBJECTS_TYPE_CHECKBOX,
	}
}

func NewFormObjectBarcode(name, label string) FormObject {
	return FormObject{
		Name:  name,
		Label: label,
		Value: "barcode",
		Type:  FORM_OBJECTS_TYPE_INBUILT,
	}
}
func NewFormObjectQrcode(name, label string) FormObject {
	return FormObject{
		Name:  name,
		Label: label,
		Value: "qrcode",
		Type:  FORM_OBJECTS_TYPE_INBUILT,
	}
}
func NewFormObjectSubmit(name, label string) FormObject {
	return FormObject{
		Name:  name,
		Label: label,
		Type:  FORM_OBJECTS_TYPE_SUBMIT,
	}
}
func NewFormObjectInputWithValue(name, label, value string) FormObject {
	return FormObject{
		Name:  name,
		Label: label,
		Value: value,
		Type:  FORM_OBJECTS_TYPE_TEXT,
	}
}

func NewFormObjectRadioInput(name, label string, options []FormObjectOption) FormObject {
	return FormObject{
		Name:    name,
		Label:   label,
		Type:    FORM_OBJECTS_TYPE_RADIO,
		Options: options,
	}
}
func NewFormObjectSelect(name, label string, options []FormObjectOption) FormObject {
	return FormObject{
		Name:    name,
		Label:   label,
		Type:    FORM_OBJECTS_TYPE_SELECT,
		Options: options,
	}
}
func NewForm(formObject ...FormObject) []FormObject {
	var form []FormObject
	form = append(form, formObject...)
	return form
}

func NewPhoto(chatID int64, file RequestFileData) PhotoConfig {
	return PhotoConfig{
		BaseFile: BaseFile{
			BaseChat: BaseChat{ChatID: chatID},
			File:     file,
		},
	}
}

func NewVideo(chatID int64, file RequestFileData) VideoConfig {
	return VideoConfig{
		BaseFile: BaseFile{
			BaseChat: BaseChat{ChatID: chatID},
			File:     file,
		},
	}
}

func NewVoice(chatID int64, file RequestFileData) VoiceConfig {
	return VoiceConfig{
		BaseFile: BaseFile{
			BaseChat: BaseChat{ChatID: chatID},
			File:     file,
		},
	}
}

func NewAudio(chatID int64, file RequestFileData) AudioConfig {
	return AudioConfig{
		BaseFile: BaseFile{
			BaseChat: BaseChat{ChatID: chatID},
			File:     file,
		},
	}
}

func NewFile(chatID int64, file RequestFileData) FileConfig {
	return FileConfig{
		BaseFile: BaseFile{
			BaseChat: BaseChat{ChatID: chatID},
			File:     file,
		},
	}
}

func NewAnswerCallback(chatID int64, callbackId string, text string, showAlert bool) CallbackAnswerConfig {
	return CallbackAnswerConfig{
		BaseChat:   BaseChat{ChatID: chatID},
		CallbackId: callbackId,
		Text:       text,
		ShowAlert:  showAlert,
	}
}

func (message *Message) UnmarshalJson(data []byte) error {
	err := json.Unmarshal(data, &message)
	if message.ChatID == 0 {
		var raw map[string]interface{}
		if err := json.Unmarshal(data, &raw); err != nil {
			return err
		}
		if chatIDStr, ok := raw["chat_id"].(string); ok {
			chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
			if err != nil {
				return err
			}
			message.ChatID = chatID
		}
	}
	if message.Data != "" {
		return nil
	}
	return err
}
