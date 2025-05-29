package models

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type (
	User struct {
		Id        int64  `json:"id,omitempty"`
		UUId      string `json:"uu_id,omitempty"`
		Username  string `json:"user,omitempty"`
		Name      string `json:"name,omitempty"`
		IsDeleted bool   `json:"is_deleted,omitempty"`
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

	Location struct {
		Lat  string `json:"lat"`
		Long string `json:"long"`
		Desc string `json:"desc,omitempty"`
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
		Path        string    `json:"path,omitempty"`
		Screenshots ImageUrls `json:"image_urls,omitempty"`
	}

	FileDta struct {
		File
		Description string `json:"desc"`
	}

	FormObject struct {
		Name    string             `json:"name,omitempty"`
		Type    string             `json:"type,omitempty"`
		Label   string             `json:"label,omitempty"`
		Value   string             `json:"value,omitempty"`
		Options []FormObjectOption `json:"options,omitempty"`
	}

	FormObjectOption map[string]string

	PaymentInfo struct {
		RefId     string `json:"ref_id"`
		MessageId string `json:"message_id"`
		Status    string `json:"status"`
	}

	FormData struct {
		MessageID  int64             `json:"message_id"`
		CallbackID string            `json:"callback_id"`
		RowData    string            `json:"data"`
		Data       map[string]string `json:"-"`
	}

	CallbackQuery struct {
		MessageID  int64               `json:"message_id"`
		Data       string              `json:"data"`
		QueryActin CallbackQueryAction `json:"-"`
		CallbackId string              `json:"callback_id"`
	}

	CallbackQueryAction struct {
		StatePath string            `json:"state_path"`
		Params    map[string]string `json:"params"`
	}

	Message struct {
		ChatID        int64         `json:"chat_id"`
		MessageID     int64         `json:"id"`
		Text          string        `json:"text"`
		Data          string        `json:"data"`
		From          User          `json:"from"`
		Type          string        `json:"type"`
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
		Text         string `json:"text"`
		CallbackData string `json:"cb_data,omitempty"`
		URL          string `json:"url,omitempty"`
		OpenIn       string `json:"open_in,omitempty"`
		Amount       int    `json:"amount,omitempty"`
		Currency     string `json:"currency,omitempty"`
		RefId        string `json:"ref_id,omitempty"`
		Description  string `json:"desc,omitempty"`
	}

	InlineKeyboardMarkup [][]InlineKeyboardButton

	APIResponse struct {
		Error     string `json:"error"`
		TraceId   string `json:"trace_id"`
		MessageId int64  `json:"id"`
	}

	Error struct {
		Message string
	}
)

func (e Error) Error() string {
	return e.Message
}

func (m *Message) UnmarshalJSON(data []byte) error {
	type Alias Message
	aux := &struct {
		ChatID interface{} `json:"chat_id"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	switch v := aux.ChatID.(type) {
	case float64:
		m.ChatID = int64(v)
	case string:
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid chat_id string: %w", err)
		}
		m.ChatID = id
	default:
		m.ChatID = 0
	}

	return nil
}

func (ikm InlineKeyboardMarkup) AddRow(row []InlineKeyboardButton) InlineKeyboardMarkup {
	return append(ikm, row)
}

func (ikm InlineKeyboardMarkup) AddButtonToEndRow(btn InlineKeyboardButton) InlineKeyboardMarkup {
	if len(ikm) == 0 {
		ikm = append(ikm, []InlineKeyboardButton{})
	}
	ikm[len(ikm)-1] = append(ikm[len(ikm)-1], btn)
	return ikm
}
