package gapBotApi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type (
	Handler func(ctx *Ctx) error
	Ctx     struct {
		bot      *BotAPI
		Endpoint string
		Message  *Message
		Params   map[string]interface{}
		context.Context
		HandlerIndex uint
		UserState    UserState
	}
	State struct {
		Endpoint string
		Message  *Message
		Params   map[string]interface{}
	}
	UserState struct {
		Stack []State
		Next  *State
	}
)

func (ctx *Ctx) Handlers() []Handler {
	handlers := make([]Handler, 0)
	if ctx.Message == nil {
		return handlers
	}

	var endpoint string
	switch ctx.Message.Type {
	case MESSAGE_TYPE_TRIGGER_BUTTON:
		endpoint = ctx.Message.CallbackQuery.QueryActin.StatePath
	default:
		endpoint = ctx.Message.Text
	}
	handlers = append(handlers, ctx.bot.Handlers[endpoint]...)
	var userState UserState
	var ok bool
	if userState, ok = ctx.bot.userStats[ctx.Message.From.Id]; !ok {
		userState = UserState{
			Stack: make([]State, 0),
			Next:  nil,
		}
	}
	var lastState = State{}
	if len(userState.Stack) > 0 {
		lastState = userState.Stack[len(userState.Stack)-1]
	}

	if endpoint != "/back" && lastState.Endpoint != endpoint {
		userState.Stack = append(userState.Stack, State{
			Endpoint: endpoint,
			Message:  ctx.Message,
			Params:   ctx.Params,
		})
	}

	if len(handlers) == 0 && userState.Next != nil {
		handlers = ctx.bot.Handlers[userState.Next.Endpoint]
		if userState.Next.Params != nil {
			for k, v := range userState.Next.Params {
				ctx.Params[k] = v
			}
		}
		userState.Next = nil
	}
	ctx.Endpoint = endpoint
	ctx.UserState = userState
	ctx.bot.userStats[ctx.Message.From.Id] = userState
	return handlers
}

func (ctx *Ctx) Middlewares() []Handler {
	return ctx.bot.Middlewares
}
func (ctx *Ctx) ResetUserStack() {
	_, ok := ctx.bot.userStats[ctx.Message.From.Id]
	if ok {
		ctx.bot.userStats[ctx.Message.From.Id] = UserState{
			Stack: make([]State, 0),
			Next:  nil,
		}
	}
}
func (ctx *Ctx) Back() error {
	if len(ctx.UserState.Stack) > 0 {
		ctx.CleanState()
		if len(ctx.UserState.Stack) > 0 {
			previousCtx := &Ctx{
				bot:     ctx.bot,
				Message: ctx.UserState.Stack[len(ctx.UserState.Stack)-1].Message,
				Params:  ctx.UserState.Stack[len(ctx.UserState.Stack)-1].Params,
				Context: ctx.Context,
			}
			return previousCtx.Next()
		}
	}
	return nil
}
func (ctx *Ctx) CleanState() {
	if userState, ok := ctx.bot.userStats[ctx.Message.From.Id]; ok && len(userState.Stack) > 0 {
		userState.Stack = userState.Stack[:len(userState.Stack)-1]
		ctx.bot.userStats[ctx.Message.From.Id] = userState
		ctx.UserState = userState
	}
}

func (ctx *Ctx) Next() error {
	handlers := ctx.Handlers()
	if len(handlers) == 0 && ctx.bot.DefaultHandler != nil {
		handlers = append(handlers, *ctx.bot.DefaultHandler)
	}

	if len(handlers) > 0 {
		handlers = append(ctx.Middlewares(), handlers...)
		if ctx.HandlerIndex <= uint(len(handlers)) {
			handler := handlers[ctx.HandlerIndex]
			if handler != nil {
				ctx.HandlerIndex++
				return handler(ctx)
			}
		}
	}
	return nil
}

func (ctx *Ctx) Unmarshal(update []byte) error {
	err := ctx.Message.UnmarshalJson(update)
	if err != nil {
		return err
	}

	switch ctx.Message.Type {
	case MESSAGE_TYPE_TEXT:
		ctx.Message.Text = ctx.Message.Data
	case MESSAGE_TYPE_IMAGE:
		err = json.Unmarshal([]byte(ctx.Message.Data), &ctx.Message.Photo)
		if err != nil {
			return fmt.Errorf("unmarshal image: %w", err)
		}
	case MESSAGE_TYPE_VIDEO:
		err = json.Unmarshal([]byte(ctx.Message.Data), &ctx.Message.Video)
		if err != nil {
			return fmt.Errorf("unmarshal video: %w", err)
		}
	case MESSAGE_TYPE_FILE:
		err = json.Unmarshal([]byte(ctx.Message.Data), &ctx.Message.File)
		if err != nil {
			return fmt.Errorf("unmarshal file: %w", err)
		}
	case MESSAGE_TYPE_AUDIO:
		err = json.Unmarshal([]byte(ctx.Message.Data), &ctx.Message.Audio)
		if err != nil {
			return fmt.Errorf("unmarshal audio: %w", err)
		}
	case MESSAGE_TYPE_VOICE:
		err = json.Unmarshal([]byte(ctx.Message.Data), &ctx.Message.Voice)
		if err != nil {
			return fmt.Errorf("unmarshal voice: %w", err)
		}
	case MESSAGE_TYPE_LOCATION:
		err = json.Unmarshal([]byte(ctx.Message.Data), &ctx.Message.Location)
		if err != nil {
			return fmt.Errorf("unmarshal location: %w", err)
		}
	case MESSAGE_TYPE_CONTACT:
		err = json.Unmarshal([]byte(ctx.Message.Data), &ctx.Message.Contact)
		if err != nil {
			return fmt.Errorf("unmarshal contact: %w", err)
		}
	case MESSAGE_TYPE_PAY_CALLBACK:
		err = json.Unmarshal([]byte(ctx.Message.Data), &ctx.Message.PaymentInfo)
		if err != nil {
			return fmt.Errorf("unmarshal payment info: %w", err)
		}
	case MESSAGE_TYPE_SUBMITFORM:
		err = json.Unmarshal([]byte(ctx.Message.Data), &ctx.Message.FormData)
		if err != nil {
			return fmt.Errorf("unmarshal form data: %w", err)
		}

		strVal := ctx.Message.FormData.RowData
		strVals := strings.SplitN(strVal, "?", 2)
		if len(strVals) > 1 {
			strVal = strVals[1]
		}

		values, err := url.ParseQuery(strVal)
		if err != nil {
			return fmt.Errorf("parse form data: %w", err)
		}

		result := make(map[string]string)
		for key, val := range values {
			if len(val) > 0 {
				result[key] = val[0]
			}
		}
		ctx.Message.FormData.Data = result

	case MESSAGE_TYPE_TRIGGER_BUTTON:
		err = json.Unmarshal([]byte(ctx.Message.Data), &ctx.Message.CallbackQuery)
		if err != nil {
			return fmt.Errorf("unmarshal callback query: %w", err)
		}

		err = json.Unmarshal([]byte(ctx.Message.CallbackQuery.Data), &ctx.Message.CallbackQuery.QueryActin)
		if err != nil {
			return fmt.Errorf("unmarshal callback query action: %w", err)
		}
		ctx.Message.MessageID = ctx.Message.CallbackQuery.MessageID

	case MESSAGE_TYPE_JOIN:
		ctx.Message.Text = "/start"
	case MESSAGE_TYPE_LEAVE:
		ctx.Message.Text = "/leave"
	}
	return nil
}

func (ctx *Ctx) SetNextStat(state State) {
	if userState, ok := ctx.bot.userStats[ctx.Message.From.Id]; ok {
		userState.Next = &state
		ctx.bot.userStats[ctx.Message.From.Id] = userState
	}
}

func (ctx *Ctx) Bot() *BotAPI {
	return ctx.bot
}

func (ctx *Ctx) WithParam(key string, val interface{}) *Ctx {
	ctx.Params[key] = val
	return ctx
}

func (ctx *Ctx) GetParam(key string) interface{} {
	if val, ok := ctx.Params[key]; ok {
		return val
	}
	return nil
}

func (ctx *Ctx) WithContext(newCtx context.Context) {
	ctx.Context = newCtx
}

func (ctx *Ctx) GetContext() context.Context {
	return ctx.Context
}
