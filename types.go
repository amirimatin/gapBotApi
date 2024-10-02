package gapBotApi

type (
	MessageHandlers    map[string]MessageHandlerFunc
	CallbackHandlers   map[string]CallbackHandlerFunc
	MiddlewareHandlers []MiddlewareHandlerFunc

	MessageHandlerFunc    func(botApi *BotAPI, message *Message) error
	CallbackHandlerFunc   func(botApi *BotAPI, callback *CallbackQuery, params map[string]string) error
	MiddlewareHandlerFunc func(botApi *BotAPI, message *Message)

	CallbackQuery struct {
		ChatId     int64               `json:"-"`
		MessageID  int64               `json:"message_id"`
		UserId     int64               `json:"user_id"`
		Data       string              `json:"data"`
		QueryActin CallbackQueryAction `json:"-"`
		CallbackId string              `json:"callback_id"`
	}
	CallbackQueryAction struct {
		StatePath string            `json:"state_path"`
		Params    map[string]string `json:"params"`
	}
	User struct {
		Id        int64  `json:"id,omitempty"`
		UUId      string `json:"uu_id,omitempty"`
		Username  string `json:"user,omitempty"`
		Name      string `json:"name,omitempty"`
		IsDeleted bool   `json:"is_deleted,omitempty"`
		//Avatar    string `json:"avatar"`
	}

	ImageUrls struct {
		Url64  string `json:"64,omitempty"`
		Url128 string `json:"128,omitempty"`
		Url256 string `json:"256,omitempty"`
		Url512 string `json:"512,omitempty"`
	}
	Contact struct {
		Id          int64  `json:"id"`
		PhoneNumber string `json:"phone"`
		Name        string `json:"name,omitempty"`
	}

	FormObject struct {
		Name    string             `json:"name,omitempty"`
		Type    FORM_OBJECTS_TYPE  `json:"type,omitempty"`
		Label   string             `json:"label,omitempty"`
		Value   string             `json:"value,omitempty"`
		Options []FormObjectOption `json:"options,omitempty"`
	}
	FormObjectOption map[string]string

	Location struct {
		Lat  string `json:"lat"`
		Long string `json:"long"`
		Desc string `json:"desc,omitempty"`
	}

	PaymentInfo struct {
		RefId     string `json:"ref_id"`
		MessageId string `json:"message_id"`
		Status    string `json:"status"`
	}

	APIResponse struct {
		Error     string `json:"error"`
		TraceId   string `json:"trace_id"`
		MessageId int64  `json:"id"`
	}

	File struct {
		Id          int64     `json:"id,omitempty"`
		SID         string    `json:"SID,omitempty"`
		RoundVideo  bool      `json:"RoundVideo,omitempty"`
		Extension   string    `json:"extension,omitempty"`
		Filename    string    `json:"filename,omitempty"`
		Filesize    int64     `json:"filesize,omitempty"`
		Type        string    `json:"type,omitempty"`
		Width       int64     `json:"width,omitempty"`
		Height      int64     `json:"height,omitempty"`
		Duration    float64   `json:"duration,omitempty"`
		Desc        string    `json:"desc,omitempty"`
		Path        string    `json:"path,omitempty"`
		Screenshots ImageUrls `json:"image_urls,omitempty"`
	}
	FormData struct {
		MessageID  int64             `json:"message_id"`
		CallbackID string            `json:"callback_id"`
		RowData    string            `json:"data"`
		Data       map[string]string `json:"-"`
	}
	UploadResponse struct {
		APIResponse `json:"-"`
		File        `json:"-"`
	}

	// Error is an error containing extra information returned by the Telegram API.
	Error struct {
		Message string
	}

	// Message represents a messageHandler.
	Message struct {
		ChatID        int64         `json:"chat_id"`
		MessageID     int64         `json:"id"`
		Text          string        `json:"text"`
		Data          string        `json:"data"`
		From          User          `json:"from"`
		Type          MESSAGE_TYPE  `json:"type"`
		Photo         File          `json:"photo,omitempty"`
		Video         File          `json:"video,omitempty"`
		Voice         File          `json:"voice,omitempty"`
		Audio         File          `json:"audio,omitempty"`
		File          File          `json:"file,omitempty"`
		PaymentInfo   PaymentInfo   `json:"payment_info,omitempty"`
		CallbackQuery CallbackQuery `json:"callback,omitempty"`
		Contact       Contact       `json:"contact,omitempty"`
		Location      Location      `json:"location,omitempty"`
		FormData      FormData      `json:"form_data,omitempty"`
	}
	ReplyKeyboardButton map[string]string
	ReplyKeyboardMarkup struct {
		Keyboard [][]ReplyKeyboardButton `json:"keyboard"`
	}

	InlineKeyboardButton struct {
		Text         string                     `json:"text"`
		CallbackData string                     `json:"cb_data,omitempty"`
		URL          string                     `json:"url,omitempty"`
		OpenIn       INLINE_KEYBOARD_URL_OPENIN `json:"open_in,omitempty"`
		Amount       int                        `json:"amount,omitempty"`
		Currency     INLINE_KEYBOARD_CURRENCY   `json:"currency,omitempty"`
		RefId        string                     `json:"ref_id,omitempty"`
		Description  string                     `json:"desc,omitempty"`
	}

	InlineKeyboardMarkup [][]InlineKeyboardButton
)

func (e Error) Error() string {
	return e.Message
}
