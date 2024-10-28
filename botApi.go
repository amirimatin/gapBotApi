package gapBotApi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"io"
	"log"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"strconv"
	"time"
)

// BotAPI structure with Resty client
type BotAPI struct {
	Token           string               `json:"token"`
	Debug           bool                 `json:"debug"`
	Client          *resty.Client        `json:"-"`
	Handlers        map[string][]Handler `json:"-"`
	Middlewares     []Handler            `json:"-"`
	DefaultHandler  *Handler             `json:"-"`
	userStats       map[int64]UserState
	shutdownChannel chan interface{}
	apiEndpoint     string
}

// NewBotAPI creates a new BotAPI instance.
func NewBotAPI(token string) (*BotAPI, error) {
	return NewBotAPIWithClient(token, APIEndpoint, resty.New().SetTimeout(30*time.Second))
}

func NewBotAPIWithClient(token, apiEndpoint string, client *resty.Client) (*BotAPI, error) {
	client.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	client.SetHeader("token", token)

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

// GetHandlers retrieves handlers for a specific endpoint.
func (bot *BotAPI) GetHandlers(endpoint string) []Handler {
	return bot.Handlers[endpoint]
}

func (bot *BotAPI) Use(handler ...Handler) {
	bot.Middlewares = append(bot.Middlewares, handler...)
}

func (bot *BotAPI) Handle(endpoint string, handler ...Handler) {
	bot.Handlers[endpoint] = append(bot.Handlers[endpoint], handler...)
}

// MakeRequest handles HTTP requests to the specified endpoint with Resty.
func (bot *BotAPI) MakeRequest(endpoint string, params Params) (*APIResponse, error) {
	if bot.Debug {
		log.Printf("Endpoint: %s, params: %v\n", endpoint, params)
	}

	method := fmt.Sprintf(bot.apiEndpoint, endpoint)
	values := buildParams(params)

	var apiResp APIResponse

	_, err := bot.Client.R().
		SetBody(values.Encode()).
		SetResult(&apiResp).
		Post(method)

	if err != nil {
		return nil, err
	}
	if apiResp.Error != "" {
		return &apiResp, &Error{
			Message: apiResp.Error,
		}
	}

	return &apiResp, nil
}

func buildParams(in Params) url.Values {
	out := url.Values{}
	for key, value := range in {
		out.Set(key, value)
	}
	return out
}

// Request sends a request and handles file uploads if needed.
func (bot *BotAPI) Request(c Chattable) (*APIResponse, error) {
	params, err := c.params()
	if err != nil {
		return nil, err
	}

	if t, ok := c.(Fileable); ok {
		file := t.file()
		if hasFileNeedingUpload(file) {
			uFile, err := bot.UploadFile(params, file)
			if err != nil {
				return nil, err
			}
			var mFile FileDta
			mFile.File = *uFile
			mFile.Description = params["desc"]
			stringFileData, err := json.Marshal(mFile)
			if err != nil {
				return nil, err
			}
			params["data"] = string(stringFileData)
		} else {
			params["data"] = file.Data.SendData()
		}
	}
	var apiResp APIResponse

	resp, err := bot.Client.R().
		SetFormData(params).
		//SetResult(&apiResp).
		SetHeader("token", bot.Token).
		Post(fmt.Sprintf(bot.apiEndpoint, c.method()))

	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp.Body(), &apiResp)
	if err != nil {
		return nil, err
	}
	if apiResp.Error != "" {
		return &apiResp, &Error{
			Message: apiResp.Error,
		}
	}

	return &apiResp, nil
}

func hasFileNeedingUpload(file RequestFile) bool {
	return file.Data.NeedsUpload()
}

// UploadFile uploads files using Resty.
func (bot *BotAPI) UploadFile(params Params, file RequestFile) (*File, error) {
	w := &bytes.Buffer{}
	m := multipart.NewWriter(w)
	defer m.Close()

	for field, value := range params {
		if err := m.WriteField(field, value); err != nil {
			return nil, err
		}
	}

	if file.Data.NeedsUpload() {
		name, reader, err := file.Data.UploadData()
		if err != nil {
			return nil, err
		}

		part, err := m.CreateFormFile(file.Name, filepath.Base(name))
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(part, reader)
		if err != nil {
			return nil, err
		}

		err = m.Close()
		if err != nil {
			return nil, err
		}
		var mFile File
		resp, err := bot.Client.R().
			SetHeader("Content-Type", m.FormDataContentType()).
			SetBody(w).
			SetResult(&mFile).
			Post(fmt.Sprintf(bot.apiEndpoint, "upload"))
		if err != nil {
			return nil, err
		}
		fmt.Println(string(resp.Body()))

		if mFile.SID == "" {
			var apiResp APIResponse
			err := json.Unmarshal(resp.Body(), &apiResp)
			if err != nil {
				return nil, err
			}
			if apiResp.Error != "" {
				return nil, &Error{
					Message: apiResp.Error,
				}
			}
		}

		return &mFile, nil
	}

	return nil, errors.New("no file to upload")
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

func (bot *BotAPI) HandleUpdates(update []byte) error {
	ctx := Ctx{
		bot:          bot,
		Message:      &Message{},
		Context:      context.Background(),
		HandlerIndex: 0,
		Params:       make(map[string]interface{}),
	}
	err := ctx.Unmarshal(update)
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

	if err := app.Listen(fmt.Sprintf(":%v", port)); err != nil {
		panic(err)
	}
}
