package gapBotApi

import (
	"io"
	"os"
)

// Telegram constants
const (
	// APIEndpoint is the Endpoint for all API methods,
	// with formatting for Sprintf.
	APIEndpoint = "https://api.gap.im/%s"
)

type MESSAGE_TYPE string

const (
	MESSAGE_TYPE_JOIN             MESSAGE_TYPE = "join"
	MESSAGE_TYPE_LEAVE            MESSAGE_TYPE = "leave"
	MESSAGE_TYPE_TEXT             MESSAGE_TYPE = "text"
	MESSAGE_TYPE_IMAGE            MESSAGE_TYPE = "image"
	MESSAGE_TYPE_AUDIO            MESSAGE_TYPE = "audio"
	MESSAGE_TYPE_VIDEO            MESSAGE_TYPE = "video"
	MESSAGE_TYPE_VOICE            MESSAGE_TYPE = "voice"
	MESSAGE_TYPE_FILE             MESSAGE_TYPE = "file"
	MESSAGE_TYPE_CONTACT          MESSAGE_TYPE = "contact"
	MESSAGE_TYPE_LOCATION         MESSAGE_TYPE = "location"
	MESSAGE_TYPE_SUBMITFORM       MESSAGE_TYPE = "submitForm"
	MESSAGE_TYPE_TRIGGER_BUTTON   MESSAGE_TYPE = "triggerButton"
	MESSAGE_TYPE_PAY_CALLBACK     MESSAGE_TYPE = "paycallback"
	MESSAGE_TYPE_INVOICE_CALLBACK MESSAGE_TYPE = "invoicecallback"
)

// Constant values for ChatActions
const (
	ChatTyping          = "typing"
	ChatUploadPhoto     = "upload_photo"
	ChatRecordVideo     = "record_video"
	ChatUploadVideo     = "upload_video"
	ChatRecordVoice     = "record_voice"
	ChatUploadVoice     = "upload_voice"
	ChatUploadDocument  = "upload_document"
	ChatChooseSticker   = "choose_sticker"
	ChatFindLocation    = "find_location"
	ChatRecordVideoNote = "record_video_note"
	ChatUploadVideoNote = "upload_video_note"
)

type INLINE_KEYBOARD_URL_OPENIN string

const (
	INLINE_KEYBOARD_URL_OPENIN_BROWSER             INLINE_KEYBOARD_URL_OPENIN = "browser"
	INLINE_KEYBOARD_URL_OPENIN_INLINE_BROWSER      INLINE_KEYBOARD_URL_OPENIN = "inline_browser"
	INLINE_KEYBOARD_URL_OPENIN_WEBVIEW             INLINE_KEYBOARD_URL_OPENIN = "webview"
	INLINE_KEYBOARD_URL_OPENIN_WEBVIEW_FULL        INLINE_KEYBOARD_URL_OPENIN = "webview_full"
	INLINE_KEYBOARD_URL_OPENIN_WEBVIEW_WITH_HEADER INLINE_KEYBOARD_URL_OPENIN = "webview_with_header"
)

type INLINE_KEYBOARD_CURRENCY string

const (
	INLINE_KEYBOARD_CURRENCY_IRR   INLINE_KEYBOARD_CURRENCY = "IRR"
	INLINE_KEYBOARD_CURRENCY_GAPCY INLINE_KEYBOARD_CURRENCY = "coin"
)

type FORM_OBJECTS_TYPE string

const (
	FORM_OBJECTS_TYPE_TEXT     FORM_OBJECTS_TYPE = "text"
	FORM_OBJECTS_TYPE_RADIO    FORM_OBJECTS_TYPE = "radio"
	FORM_OBJECTS_TYPE_SELECT   FORM_OBJECTS_TYPE = "select"
	FORM_OBJECTS_TYPE_TEXTAREA FORM_OBJECTS_TYPE = "textarea"
	FORM_OBJECTS_TYPE_INBUILT  FORM_OBJECTS_TYPE = "inbuilt"
	FORM_OBJECTS_TYPE_CHECKBOX FORM_OBJECTS_TYPE = "checkbox"
	FORM_OBJECTS_TYPE_SUBMIT   FORM_OBJECTS_TYPE = "submit"
)

type Chattable interface {
	params() (Params, error)
	method() string
}

type Fileable interface {
	Chattable
	file() RequestFile
}

type RequestFile struct {
	// The file field name.
	Name string
	Type MESSAGE_TYPE
	// The file data to include.
	Data RequestFileData
}

// FileReader contains information about a reader to upload as a File.
type FileReader struct {
	Name   string
	Reader io.Reader
}

func (fr FileReader) NeedsUpload() bool {
	return true
}

func (fr FileReader) UploadData() (string, io.Reader, error) {
	return fr.Name, fr.Reader, nil
}

func (fr FileReader) SendData() string {
	panic("FileReader must be uploaded")
}

// FilePath is a path to a local file.
type FilePath string

func (fp FilePath) NeedsUpload() bool {
	return true
}

func (fp FilePath) UploadData() (string, io.Reader, error) {
	fileHandle, err := os.Open(string(fp))
	if err != nil {
		return "", nil, err
	}

	name := fileHandle.Name()
	return name, fileHandle, err
}
func (fp FilePath) SendData() string {
	panic("FilePath must be uploaded")
}

type RequestFileData interface {
	// NeedsUpload shows if the file needs to be uploaded.
	NeedsUpload() bool

	// UploadData gets the file name and an `io.Reader` for the file to be uploaded. This
	// must only be called when the file needs to be uploaded.
	UploadData() (string, io.Reader, error)
	// SendData gets the file data to send when a file does not need to be uploaded. This
	// must only be called when the file does not need to be uploaded.
	SendData() string
}

// BaseChat is base type for all chat config types.
type BaseChat struct {
	ChatID               int64 // required
	ReplyKeyboardMarkup  interface{}
	InlineKeyboardMarkup InlineKeyboardMarkup
}

func (chat *BaseChat) params() (Params, error) {
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

type MessageConfig struct {
	BaseChat
	Type MESSAGE_TYPE `json:"type"`
	Text string       `json:"data"`
	Form []FormObject `json:"form"`
}

func (config MessageConfig) params() (Params, error) {
	params, err := config.BaseChat.params()
	if err != nil {
		return params, err
	}
	params.AddNonEmpty("data", config.Text)
	params.AddNonEmpty("type", "text")
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
	Type      MESSAGE_TYPE `json:"type"`
	Text      string       `json:"data"`
	MessageId int64        `json:"message_id"`
}

func (config UpdateMessageConfig) params() (Params, error) {
	params, err := config.BaseChat.params()
	if err != nil {
		return params, err
	}
	params.AddNonEmpty("data", config.Text)
	params.AddNonZero64("message_id", config.MessageId)
	return params, err
}

func (config UpdateMessageConfig) method() string {
	return "editMessage"
}

type DeleteMessageConfig struct {
	BaseChat
	Type      MESSAGE_TYPE `json:"type"`
	Text      string       `json:"data"`
	MessageId int64        `json:"message_id"`
}

func (config DeleteMessageConfig) params() (Params, error) {
	params, err := config.BaseChat.params()
	if err != nil {
		return params, err
	}
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
	params, err := config.BaseChat.params()
	if err != nil {
		return params, err
	}
	params.AddNonEmpty("text", config.Text)
	params.AddBool("show_alert", config.ShowAlert)
	params.AddNonEmpty("callback_id", config.CallbackId)
	return params, err
}

func (config CallbackAnswerConfig) method() string {
	return "answerCallback"
}

// BaseFile is a base type for all file config types.
type BaseFile struct {
	BaseChat
	File RequestFileData
}

func (file BaseFile) params() (Params, error) {
	return file.BaseChat.params()
}

type PhotoConfig struct {
	BaseFile
	Description string
}

func (config PhotoConfig) params() (Params, error) {
	params, err := config.BaseFile.params()
	if err != nil {
		return params, err
	}

	params.AddNonEmpty("desc", config.Description)
	params.AddNonEmpty("type", "image")
	return params, err
}

func (config PhotoConfig) method() string {
	return "sendMessage"
}

func (config PhotoConfig) file() RequestFile {
	return RequestFile{
		Name: "image",
		Type: MESSAGE_TYPE_IMAGE,
		Data: config.File,
	}

}

type VideoConfig struct {
	BaseFile
	Description string
}

func (config VideoConfig) params() (Params, error) {
	params, err := config.BaseFile.params()
	if err != nil {
		return params, err
	}

	params.AddNonEmpty("desc", config.Description)
	params.AddNonEmpty("type", "video")
	return params, err
}

func (config VideoConfig) method() string {
	return "sendMessage"
}

func (config VideoConfig) file() RequestFile {
	return RequestFile{
		Name: "video",
		Type: MESSAGE_TYPE_VIDEO,
		Data: config.File,
	}
}

type VoiceConfig struct {
	BaseFile
	Description string
}

func (config VoiceConfig) params() (Params, error) {
	params, err := config.BaseFile.params()
	if err != nil {
		return params, err
	}

	params.AddNonEmpty("desc", config.Description)
	params.AddNonEmpty("type", "voice")
	return params, err
}

func (config VoiceConfig) method() string {
	return "sendMessage"
}

func (config VoiceConfig) file() RequestFile {
	return RequestFile{
		Name: "voice",
		Type: MESSAGE_TYPE_VOICE,
		Data: config.File,
	}
}

type AudioConfig struct {
	BaseFile
	Description string
}

func (config AudioConfig) params() (Params, error) {
	params, err := config.BaseFile.params()
	if err != nil {
		return params, err
	}

	params.AddNonEmpty("desc", config.Description)
	params.AddNonEmpty("type", "audio")
	return params, err
}

func (config AudioConfig) method() string {
	return "sendMessage"
}

func (config AudioConfig) file() RequestFile {
	return RequestFile{
		Name: "audio",
		Type: MESSAGE_TYPE_AUDIO,
		Data: config.File,
	}
}

type FileConfig struct {
	BaseFile
	Description string
}

func (config FileConfig) params() (Params, error) {
	params, err := config.BaseFile.params()
	if err != nil {
		return params, err
	}

	params.AddNonEmpty("desc", config.Description)
	params.AddNonEmpty("type", "file")
	return params, err
}

func (config FileConfig) method() string {
	return "sendMessage"
}

func (config FileConfig) file() RequestFile {
	return RequestFile{
		Name: "file",
		Type: MESSAGE_TYPE_FILE,
		Data: config.File,
	}
}
