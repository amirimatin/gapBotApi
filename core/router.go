// router.go
package gapBotApi

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/amirimatin/gapBotApi/v2/models"
)

func (ctx *Ctx) Handlers() []Handler {
	var handlers []Handler
	if ctx.Message == nil {
		return handlers
	}

	var endpoint string
	if ctx.Message.Type == MESSAGE_TYPE_TRIGGER_BUTTON {
		endpoint = ctx.Message.CallbackQuery.QueryActin.StatePath
		for k, v := range ctx.Message.CallbackQuery.QueryActin.Params {
			ctx.Params[k] = v
		}
	} else {
		endpoint = ctx.Message.Text
	}

	if query, err := parseQuery(endpoint); err == nil {
		for k, v := range query {
			ctx.WithParam(k, v)
		}
		endpoint = strings.Split(endpoint, "?")[0]
	}

	handlers = append(handlers, ctx.bot.Handlers[endpoint]...)

	userState := ctx.bot.userStats[ctx.Message.From.Id]
	if len(userState.Stack) > 0 && endpoint != "/back" {
		userState.Stack = append(userState.Stack, State{
			Endpoint: endpoint,
			Message:  ctx.Message,
			Params:   ctx.Params,
		})
	}
	ctx.UserState = userState
	ctx.Endpoint = endpoint
	ctx.bot.userStats[ctx.Message.From.Id] = userState

	if len(handlers) == 0 && userState.Next != nil {
		handlers = ctx.bot.Handlers[userState.Next.Endpoint]
		for k, v := range userState.Next.Params {
			ctx.Params[k] = v
		}
	}

	return handlers
}

func (ctx *Ctx) Next() (models.Message, error) {
	handlers := ctx.Handlers()
	if len(handlers) == 0 && ctx.bot.DefaultHandler != nil {
		handlers = append(handlers, ctx.bot.DefaultHandler)
	}

	if len(handlers) > 0 {
		handlers = append(ctx.bot.Middlewares, handlers...)
		if ctx.HandlerIndex < uint(len(handlers)) {
			handler := handlers[ctx.HandlerIndex]
			ctx.HandlerIndex++
			return handler(ctx)
		}
	}

	return models.Message{}, nil
}

func (ctx *Ctx) Back() (models.Message, error) {
	if len(ctx.UserState.Stack) > 0 {
		ctx.CleanState()
		if len(ctx.UserState.Stack) > 0 {
			prev := ctx.UserState.Stack[len(ctx.UserState.Stack)-1]
			return (&Ctx{
				bot:     ctx.bot,
				Message: prev.Message,
				Params:  prev.Params,
				Context: ctx.Context,
			}).Next()
		}
	}
	return models.Message{}, nil
}

func (ctx *Ctx) CleanState() {
	if userState, ok := ctx.bot.userStats[ctx.Message.From.Id]; ok && len(userState.Stack) > 0 {
		userState.Stack = userState.Stack[:len(userState.Stack)-1]
		ctx.bot.userStats[ctx.Message.From.Id] = userState
		ctx.UserState = userState
	}
}

func parseQuery(queryURL string) (map[string]interface{}, error) {
	u, err := url.Parse(queryURL)
	if err != nil {
		return nil, err
	}
	result := make(map[string]interface{})
	for k, v := range u.Query() {
		if len(v) == 1 {
			result[k] = v[0]
		} else {
			result[k] = v
		}
	}
	return result, nil
}

func (ctx *Ctx) Unmarshal(update []byte) error {
	err := ctx.Message.UnmarshalJSON(update)
	if err != nil {
		return err
	}

	switch ctx.Message.Type {
	case MESSAGE_TYPE_TEXT:
		ctx.Message.Text = ctx.Message.Data
	case MESSAGE_TYPE_IMAGE:
		err = json.Unmarshal([]byte(ctx.Message.Data), &ctx.Message.Photo)
	case MESSAGE_TYPE_VIDEO:
		err = json.Unmarshal([]byte(ctx.Message.Data), &ctx.Message.Video)
	case MESSAGE_TYPE_FILE:
		err = json.Unmarshal([]byte(ctx.Message.Data), &ctx.Message.File)
	case MESSAGE_TYPE_AUDIO:
		err = json.Unmarshal([]byte(ctx.Message.Data), &ctx.Message.Audio)
	case MESSAGE_TYPE_VOICE:
		err = json.Unmarshal([]byte(ctx.Message.Data), &ctx.Message.Voice)
	case MESSAGE_TYPE_LOCATION:
		err = json.Unmarshal([]byte(ctx.Message.Data), &ctx.Message.Location)
	case MESSAGE_TYPE_CONTACT:
		err = json.Unmarshal([]byte(ctx.Message.Data), &ctx.Message.Contact)
	case MESSAGE_TYPE_PAY_CALLBACK:
		err = json.Unmarshal([]byte(ctx.Message.Data), &ctx.Message.PaymentInfo)
	case MESSAGE_TYPE_SUBMITFORM:
		err = json.Unmarshal([]byte(ctx.Message.Data), &ctx.Message.FormData)
		if err == nil {
			str := ctx.Message.FormData.RowData
			if parts := strings.SplitN(str, "?", 2); len(parts) > 1 {
				values, err := url.ParseQuery(parts[1])
				if err == nil {
					data := make(map[string]string)
					for k, v := range values {
						if len(v) > 0 {
							data[k] = v[0]
						}
					}
					ctx.Message.FormData.Data = data
				}
			}
		}
	case MESSAGE_TYPE_TRIGGER_BUTTON:
		err = json.Unmarshal([]byte(ctx.Message.Data), &ctx.Message.CallbackQuery)
		if err == nil {
			err = json.Unmarshal([]byte(ctx.Message.CallbackQuery.Data), &ctx.Message.CallbackQuery.QueryActin)
			ctx.Message.MessageID = ctx.Message.CallbackQuery.MessageID
		}
	case MESSAGE_TYPE_JOIN:
		ctx.Message.Text = "/start"
	case MESSAGE_TYPE_LEAVE:
		ctx.Message.Text = "/leave"
	}
	return err
}
