package gapBotApi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// HTTPClient is the type needed for the bot to perform HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type BotAPI struct {
	Token                string     `json:"token"`
	Debug                bool       `json:"debug"`
	Client               HTTPClient `json:"-"`
	MessageHandlers      MessageHandlers
	CallbackHandlers     CallbackHandlers
	MiddlewareHandlers   MiddlewareHandlers
	shutdownChannel      chan interface{}
	apiEndpoint          string
	DefaultTypesHandlers map[string]MessageHandlerFunc
}

// NewBotAPI creates a new BotAPI instance.
//
// It requires a token, provided by @BotFather on Telegram.
func NewBotAPI(token string) (*BotAPI, error) {
	return NewBotAPIWithClient(token, APIEndpoint, &http.Client{
		Timeout: 30 * time.Second,
	})
}

func NewBotAPIWithClient(token, apiEndpoint string, client HTTPClient) (*BotAPI, error) {
	bot := &BotAPI{
		Token:                token,
		Client:               client,
		shutdownChannel:      make(chan interface{}),
		MessageHandlers:      make(MessageHandlers),
		CallbackHandlers:     make(CallbackHandlers),
		MiddlewareHandlers:   make(MiddlewareHandlers, 0),
		apiEndpoint:          apiEndpoint,
		DefaultTypesHandlers: make(map[string]MessageHandlerFunc),
	}
	return bot, nil
}

// MakeRequest makes a request to a specific endpoint with our token.
func (bot *BotAPI) MakeRequest(endpoint string, params Params) (*APIResponse, error) {
	if bot.Debug {
		log.Printf("Endpoint: %s, params: %v\n", endpoint, params)
	}

	method := fmt.Sprintf(bot.apiEndpoint, endpoint)

	values := buildParams(params)

	req, err := http.NewRequest("POST", method, strings.NewReader(values.Encode()))
	if err != nil {
		return &APIResponse{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("token", bot.Token)
	resp, err := bot.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	bytes, err := bot.decodeAPIResponse(resp.Body, &apiResp)
	if err != nil {
		return &apiResp, err
	}

	if bot.Debug {
		log.Printf("Endpoint: %s, response: %s\n", endpoint, string(bytes))
	}

	if apiResp.Error != "" {
		return &apiResp, &Error{
			Message: apiResp.Error,
		}
	}

	return &apiResp, nil
}
func buildParams(in Params) url.Values {
	if in == nil {
		return url.Values{}
	}

	out := url.Values{}

	for key, value := range in {
		out.Set(key, value)
	}

	return out
}
func (bot *BotAPI) decodeAPIResponse(responseBody io.Reader, resp *APIResponse) ([]byte, error) {
	if !bot.Debug {
		dec := json.NewDecoder(responseBody)
		err := dec.Decode(resp)
		return nil, err
	}

	// if debug, read response body
	data, err := io.ReadAll(responseBody)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (bot *BotAPI) decodeUploadResponse(responseBody io.Reader, resp *UploadResponse) ([]byte, error) {
	if !bot.Debug {
		dec := json.NewDecoder(responseBody)
		err := dec.Decode(resp)
		return nil, err
	}

	// if debug, read response body
	data, err := io.ReadAll(responseBody)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (bot *BotAPI) Request(c Chattable) (*APIResponse, error) {
	params, err := c.params()
	if err != nil {
		return nil, err
	}

	if t, ok := c.(Fileable); ok {
		file := t.file()
		// If we have files that need to be uploaded, we should delegate the
		// request to UploadFile.
		if hasFileNeedingUpload(file) {
			uFile, err := bot.UploadFile(params, file)
			uFile.Desc = params["desc"]
			stringFileData, err := json.Marshal(uFile)
			if err != nil {
				return nil, err
			}
			params["data"] = string(stringFileData)
		} else {
			params["data"] = file.Data.SendData()
		}
	}

	return bot.MakeRequest(c.method(), params)
}
func hasFileNeedingUpload(file RequestFile) bool {
	if file.Data.NeedsUpload() {
		return true
	}
	return false
}

func (bot *BotAPI) UploadFile(params Params, file RequestFile) (*UploadResponse, error) {

	w := &bytes.Buffer{}
	m := multipart.NewWriter(w)

	defer m.Close()
	for field, value := range params {
		if err := m.WriteField(field, value); err != nil {
			return &UploadResponse{}, err
		}
	}

	if file.Data.NeedsUpload() {
		name, reader, err := file.Data.UploadData()
		if err != nil {
			return &UploadResponse{}, err
		}
		mFile, err := os.Open(name)
		if err != nil {
			fmt.Println(err)
			return &UploadResponse{}, err
		}
		defer mFile.Close()

		part, err := m.CreateFormFile(file.Name, filepath.Base(name))
		if err != nil {
			return &UploadResponse{}, err
		}

		if _, err := io.Copy(part, mFile); err != nil {
			return &UploadResponse{}, err
		}

		if closer, ok := reader.(io.ReadCloser); ok {
			if err = closer.Close(); err != nil {
				return &UploadResponse{}, err
			}
		}
		err = m.Close()
		if err != nil {
			fmt.Println(err)
			return &UploadResponse{}, err
		}
	} else {
		value := file.Data.SendData()

		if err := m.WriteField(file.Name, value); err != nil {
			return &UploadResponse{}, err
		}
	}

	if bot.Debug {
		log.Printf("Endpoint: %s, params: %v, with %d upload file\n", params)
	}

	method := fmt.Sprintf(bot.apiEndpoint, "upload")

	req, err := http.NewRequest("POST", method, w)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("token", bot.Token)
	req.Header.Set("Content-Type", m.FormDataContentType())
	resp, err := bot.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp UploadResponse
	bytes, err := bot.decodeUploadResponse(resp.Body, &apiResp)
	if err != nil {
		return &apiResp, err
	}

	if bot.Debug {
		log.Printf("Endpoint: %s, response: %s\n ", string(bytes))
	}

	if apiResp.Error != "" {
		return &apiResp, &Error{
			Message: apiResp.Error,
		}
	}

	return &apiResp, nil
}

// Send will send a Chattable item to Telegram and provides the
// returned Message.
func (bot *BotAPI) Send(c Chattable) (Message, error) {
	resp, err := bot.Request(c)
	if err != nil {
		return Message{}, err
	}
	if resp.Error != "" {
		return Message{}, errors.New(resp.Error)
	}
	msg := Message{MessageID: resp.MessageId}
	params, er := c.params()
	if er != nil {
		msg.ChatID = 0
	} else {
		i, err := strconv.ParseInt(params.GetParam("chat_id"), 10, 64)
		if err != nil {
			msg.ChatID = 0
		}
		msg.ChatID = i
	}
	return msg, err
}

func (bot *BotAPI) HandleMessage(statePath string, handler MessageHandlerFunc) {
	bot.MessageHandlers[statePath] = handler
}
func (bot *BotAPI) HandleCallback(statePath string, handler CallbackHandlerFunc) {
	bot.CallbackHandlers[statePath] = handler
}
func (bot *BotAPI) AddMiddleware(handler MiddlewareHandlerFunc) {
	bot.MiddlewareHandlers = append(bot.MiddlewareHandlers, handler)
}
func (bot *BotAPI) HandleUpdates(update []byte) (err error) {
	var message Message
	err = message.UnmarshalJson(update)
	if err != nil {
		return err
	}
	switch message.Type {
	case MESSAGE_TYPE_TEXT:
		message.Text = message.Data
	case MESSAGE_TYPE_IMAGE:
		err := json.Unmarshal([]byte(message.Data), &message.Photo)
		if err != nil {
			return err
		}
	case MESSAGE_TYPE_VIDEO:
		err := json.Unmarshal([]byte(message.Data), &message.Video)
		if err != nil {
			return err
		}
	case MESSAGE_TYPE_FILE:
		err := json.Unmarshal([]byte(message.Data), &message.File)
		if err != nil {
			return err
		}
	case MESSAGE_TYPE_AUDIO:
		err := json.Unmarshal([]byte(message.Data), &message.Audio)
		if err != nil {
			return err
		}
	case MESSAGE_TYPE_VOICE:
		err := json.Unmarshal([]byte(message.Data), &message.Voice)
		if err != nil {
			return err
		}
	case MESSAGE_TYPE_LOCATION:
		err := json.Unmarshal([]byte(message.Data), &message.Location)
		if err != nil {
			return err
		}
	case MESSAGE_TYPE_CONTACT:
		err := json.Unmarshal([]byte(message.Data), &message.Contact)
		if err != nil {
			return err
		}
	case MESSAGE_TYPE_PAY_CALLBACK:
		err := json.Unmarshal([]byte(message.Data), &message.PaymentInfo)
		if err != nil {
			return err
		}
	case MESSAGE_TYPE_SUBMITFORM:
		err := json.Unmarshal([]byte(message.Data), &message.FormData)
		if err != nil {
			return err
		}
		values, err := url.ParseQuery(message.FormData.RowData)
		if err != nil {
			return err
		}

		result := make(map[string]string)
		for key, val := range values {
			if len(val) > 0 {
				result[key] = val[0]
			}
		}
		message.FormData.Data = result
	case MESSAGE_TYPE_TRIGGER_BUTTON:
		err := json.Unmarshal([]byte(message.Data), &message.CallbackQuery)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		err = json.Unmarshal([]byte(message.CallbackQuery.Data), &message.CallbackQuery.QueryActin)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		message.MessageID = message.CallbackQuery.MessageID
		callbackData := message.CallbackQuery.QueryActin
		message.CallbackQuery.ChatId = message.ChatID
		message.CallbackQuery.UserId = message.From.Id
		bot.runMiddlewares(&message)
		handlerFunc := bot.FindCallbackHandler(callbackData.StatePath)
		if handlerFunc != nil {
			err = handlerFunc(bot, &message.CallbackQuery, callbackData.Params)
			if err != nil {
				_, err = bot.Send(NewAnswerCallback(message.ChatID, message.CallbackQuery.CallbackId, err.Error(), false))
				if err != nil {
					fmt.Printf("error in sending callback answer: %s", err.Error())
					return err
				}
			}
		} else {
			_, err := bot.Send(NewAnswerCallback(message.ChatID, message.CallbackQuery.CallbackId, "Invalid Callback Data", false))
			if err != nil {
				fmt.Printf("error in sending callback answer: %s", err.Error())
				return err
			}
		}
	case MESSAGE_TYPE_JOIN:
		message.Text = "join"
	case MESSAGE_TYPE_LEAVE:
		message.Text = "leave"
	}
	bot.runMiddlewares(&message)
	handlerFunc := bot.FindMessageHandler(message.Text)
	if handlerFunc != nil {
		return handlerFunc(bot, &message)
	}
	defaultHandler := bot.DefaultTypesHandlers[string(message.Type)]
	if defaultHandler != nil {
		err := defaultHandler(bot, &message)
		if err != nil {
			return err
		}
	}
	return nil
}
func (bot *BotAPI) FindCallbackHandler(action string) CallbackHandlerFunc {
	return bot.CallbackHandlers[action]
}

func (bot *BotAPI) FindMessageHandler(action string) MessageHandlerFunc {
	return bot.MessageHandlers[action]
}
func (bot *BotAPI) runMiddlewares(message *Message) {
	for _, handlerFunc := range bot.MiddlewareHandlers {
		handlerFunc(bot, message)
	}
}
