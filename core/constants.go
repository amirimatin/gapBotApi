package gapBotApi

const (
	// Base API endpoint
	APIEndpoint = "https://api.gap.im/%s"
)

// MESSAGE_TYPE constants represent different message types
const (
	MESSAGE_TYPE_JOIN             = "join"
	MESSAGE_TYPE_LEAVE            = "leave"
	MESSAGE_TYPE_TEXT             = "text"
	MESSAGE_TYPE_IMAGE            = "image"
	MESSAGE_TYPE_AUDIO            = "audio"
	MESSAGE_TYPE_VIDEO            = "video"
	MESSAGE_TYPE_VOICE            = "msgVoice"
	MESSAGE_TYPE_FILE             = "file"
	MESSAGE_TYPE_CONTACT          = "contact"
	MESSAGE_TYPE_LOCATION         = "location"
	MESSAGE_TYPE_SUBMITFORM       = "submitForm"
	MESSAGE_TYPE_TRIGGER_BUTTON   = "triggerButton"
	MESSAGE_TYPE_PAY_CALLBACK     = "paycallback"
	MESSAGE_TYPE_INVOICE_CALLBACK = "invoicecallback"
)

// Chat actions for bot status indication
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

// URL open-in behavior for inline keyboard buttons
const (
	INLINE_KEYBOARD_URL_OPENIN_BROWSER             = "browser"
	INLINE_KEYBOARD_URL_OPENIN_INLINE_BROWSER      = "inline_browser"
	INLINE_KEYBOARD_URL_OPENIN_WEBVIEW             = "webview"
	INLINE_KEYBOARD_URL_OPENIN_WEBVIEW_FULL        = "webview_full"
	INLINE_KEYBOARD_URL_OPENIN_WEBVIEW_WITH_HEADER = "webview_with_header"
)

// Supported currencies for inline payments
const (
	INLINE_KEYBOARD_CURRENCY_IRR   = "IRR"
	INLINE_KEYBOARD_CURRENCY_GAPCY = "coin"
)

// Supported form object types
const (
	FORM_OBJECTS_TYPE_TEXT     = "text"
	FORM_OBJECTS_TYPE_RADIO    = "radio"
	FORM_OBJECTS_TYPE_SELECT   = "select"
	FORM_OBJECTS_TYPE_TEXTAREA = "textarea"
	FORM_OBJECTS_TYPE_INBUILT  = "inbuilt"
	FORM_OBJECTS_TYPE_CHECKBOX = "checkbox"
	FORM_OBJECTS_TYPE_SUBMIT   = "submit"
)
