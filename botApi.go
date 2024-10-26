package gapBotApi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
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
	Token           string               `json:"token"`
	Debug           bool                 `json:"debug"`
	Client          HTTPClient           `json:"-"`
	Handlers        map[string][]Handler `json:"-"`
	Middlewares     []Handler            `json:"-"`
	DefaultHandler  *Handler             `json:"-"`
	userStats       map[int64]UserState
	shutdownChannel chan interface{}
	apiEndpoint     string
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
		Token:           token,
		Client:          client,
		shutdownChannel: make(chan interface{}),
		Handlers:        make(map[string][]Handler),
		Middlewares:     make([]Handler, 0),
		userStats:       make(map[int64]UserState),
		apiEndpoint:     apiEndpoint,
	}
	return bot, nil
}

// GetHandlers makes a request to a specific Endpoint with our token.
func (bot *BotAPI) GetHandlers(endpoint string) []Handler {
	handlers := bot.Handlers[endpoint]
	if handlers == nil {
		return nil
	}
	return handlers
}
func (bot *BotAPI) Use(handler ...Handler) {
	bot.Middlewares = append(bot.Middlewares, handler...)
}

func (bot *BotAPI) Handle(endpoint string, handler ...Handler) {
	bot.Handlers[endpoint] = append(bot.Handlers[endpoint], handler...)
}

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
func (bot *BotAPI) HandleUpdates(update []byte) (err error) {
	ctx := Ctx{
		bot:          bot,
		Message:      &Message{},
		Context:      context.Background(),
		HandlerIndex: 0,
	}
	err = ctx.Unmarshal(update)
	if err != nil {
		return err
	}
	return ctx.Next()
}
func (bot *BotAPI) Serve(port int, callbackEndpoint string) {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}
			return ctx.Status(code).JSON(fiber.Map{
				"status": "error",
				"error":  err,
			})
		},
		AppName: "Gap Authenticator",
	})

	bot.Handle("/back", func(ctx *Ctx) error {
		return ctx.Back()
	})

	app.Post(callbackEndpoint, func(ctx *fiber.Ctx) error {
		err := bot.HandleUpdates(ctx.Body())
		if err != nil {
			fmt.Printf("error in handle updates: %s", err.Error())
		}
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"data":   nil,
		})
	})

	err := app.Listen(fmt.Sprintf(":%v", port))
	if err != nil {
		panic(err)
	}
}
