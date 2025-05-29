package gapBotApi

import (
	"github.com/go-resty/resty/v2"
	"time"
)

type BotAPI struct {
	Token          string               `json:"token"`
	Debug          bool                 `json:"debug"`
	Client         *resty.Client        `json:"-"`
	Handlers       map[string][]Handler `json:"-"`
	Middlewares    []Handler            `json:"-"`
	DefaultHandler Handler              `json:"-"`
	userStats      map[int64]UserState
	apiEndpoint    string
}

func NewBotAPI(token string) (*BotAPI, error) {
	return NewBotAPIWithClient(token, APIEndpoint, resty.New().SetTimeout(30*time.Second))
}

func NewBotAPIWithClient(token, apiEndpoint string, client *resty.Client) (*BotAPI, error) {
	client.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	client.SetHeader("token", token)

	bot := &BotAPI{
		Token:       token,
		Client:      client,
		Handlers:    make(map[string][]Handler),
		Middlewares: make([]Handler, 0),
		userStats:   make(map[int64]UserState),
		apiEndpoint: apiEndpoint,
	}
	return bot, nil
}

func (bot *BotAPI) GetHandlers(endpoint string) []Handler {
	return bot.Handlers[endpoint]
}

func (bot *BotAPI) Use(handler ...Handler) {
	bot.Middlewares = append(bot.Middlewares, handler...)
}

func (bot *BotAPI) Handle(endpoint string, handler ...Handler) {
	bot.Handlers[endpoint] = append(bot.Handlers[endpoint], handler...)
}
