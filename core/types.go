// types.go
package gapBotApi

import (
	"github.com/amirimatin/gapBotApi/v2/models"
	. "github.com/amirimatin/gapBotApi/v2/params"
)

type BaseChat struct {
	ChatID               int64
	ReplyKeyboardMarkup  interface{}
	InlineKeyboardMarkup models.InlineKeyboardMarkup
}

type MessageConfig struct {
	BaseChat
	Type string              `json:"type"`
	Text string              `json:"data"`
	Form []models.FormObject `json:"form"`
}

func (config MessageConfig) params() (Params, error) {
	params := make(Params)
	params.AddFirstValid("chat_id", config.ChatID)
	var err error
	if config.ReplyKeyboardMarkup != nil {
		err = params.AddInterface("reply_keyboard", config.ReplyKeyboardMarkup)
	}
	if config.InlineKeyboardMarkup != nil {
		err = params.AddInterface("inline_keyboard", config.InlineKeyboardMarkup)
	}
	params.AddNonEmpty("data", config.Text)
	params.AddNonEmpty("type", config.Type)
	if len(config.Form) > 0 {
		err = params.AddInterface("form", config.Form)
	}
	return params, err
}

func (config MessageConfig) method() string {
	return "sendMessage"
}

type UpdateMessageConfig struct {
	BaseChat
	Type      string `json:"type"`
	Text      string `json:"data"`
	MessageId int64  `json:"message_id"`
}

func (config UpdateMessageConfig) params() (Params, error) {
	params, err := config.BaseChatParams()
	params.AddNonEmpty("data", config.Text)
	params.AddNonZero64("message_id", config.MessageId)
	return params, err
}

func (config UpdateMessageConfig) method() string {
	return "editMessage"
}

type DeleteMessageConfig struct {
	BaseChat
	MessageId int64 `json:"message_id"`
}

func (config DeleteMessageConfig) params() (Params, error) {
	params, err := config.BaseChatParams()
	params.AddNonZero64("message_id", config.MessageId)
	return params, err
}

func (config DeleteMessageConfig) method() string {
	return "deleteMessage"
}

type CallbackAnswerConfig struct {
	BaseChat
	CallbackId string `json:"callback_id"`
	Text       string `json:"text"`
	ShowAlert  bool   `json:"show_alert"`
}

func (config CallbackAnswerConfig) params() (Params, error) {
	params, err := config.BaseChatParams()
	params.AddNonEmpty("text", config.Text)
	params.AddBool("show_alert", config.ShowAlert)
	params.AddNonEmpty("callback_id", config.CallbackId)
	return params, err
}

func (config CallbackAnswerConfig) method() string {
	return "answerCallback"
}

type BaseFile struct {
	BaseChat
	File RequestFileData
}

func (file BaseFile) params() (Params, error) {
	return file.BaseChatParams()
}

type PhotoConfig struct {
	BaseFile
	Description string
}

func (config PhotoConfig) params() (Params, error) {
	params, err := config.BaseFile.params()
	params.AddNonEmpty("desc", config.Description)
	params.AddNonEmpty("type", MESSAGE_TYPE_IMAGE)
	return params, err
}

func (config PhotoConfig) method() string { return "sendMessage" }
func (config PhotoConfig) file() RequestFile {
	return RequestFile{Name: "image", Type: MESSAGE_TYPE_IMAGE, Data: config.File}
}

type VideoConfig struct {
	BaseFile
	Description string
}

func (config VideoConfig) params() (Params, error) {
	params, err := config.BaseFile.params()
	params.AddNonEmpty("desc", config.Description)
	params.AddNonEmpty("type", MESSAGE_TYPE_VIDEO)
	return params, err
}

func (config VideoConfig) method() string { return "sendMessage" }
func (config VideoConfig) file() RequestFile {
	return RequestFile{Name: "video", Type: MESSAGE_TYPE_VIDEO, Data: config.File}
}

type VoiceConfig struct {
	BaseFile
	Description string
}

func (config VoiceConfig) params() (Params, error) {
	params, err := config.BaseFile.params()
	params.AddNonEmpty("desc", config.Description)
	params.AddNonEmpty("type", MESSAGE_TYPE_VOICE)
	return params, err
}

func (config VoiceConfig) method() string { return "sendMessage" }
func (config VoiceConfig) file() RequestFile {
	return RequestFile{Name: "voice", Type: MESSAGE_TYPE_VOICE, Data: config.File}
}

type AudioConfig struct {
	BaseFile
	Description string
}

func (config AudioConfig) params() (Params, error) {
	params, err := config.BaseFile.params()
	params.AddNonEmpty("desc", config.Description)
	params.AddNonEmpty("type", MESSAGE_TYPE_AUDIO)
	return params, err
}

func (config AudioConfig) method() string { return "sendMessage" }
func (config AudioConfig) file() RequestFile {
	return RequestFile{Name: "audio", Type: MESSAGE_TYPE_AUDIO, Data: config.File}
}

type FileConfig struct {
	BaseFile
	Description string
}

func (config FileConfig) params() (Params, error) {
	params, err := config.BaseFile.params()
	params.AddNonEmpty("desc", config.Description)
	params.AddNonEmpty("type", MESSAGE_TYPE_FILE)
	return params, err
}

func (config FileConfig) method() string { return "sendMessage" }
func (config FileConfig) file() RequestFile {
	return RequestFile{Name: "file", Type: MESSAGE_TYPE_FILE, Data: config.File}
}

func (chat *BaseChat) BaseChatParams() (Params, error) {
	params := make(Params)
	params.AddFirstValid("chat_id", chat.ChatID)
	var err error
	if chat.ReplyKeyboardMarkup != nil {
		err = params.AddInterface("reply_keyboard", chat.ReplyKeyboardMarkup)
	}
	if chat.InlineKeyboardMarkup != nil {
		err = params.AddInterface("inline_keyboard", chat.InlineKeyboardMarkup)
	}
	return params, err
}
