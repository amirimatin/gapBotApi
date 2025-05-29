// context.go
package gapBotApi

import (
	"context"

	"github.com/amirimatin/gapBotApi/v2/models"
)

type (
	Handler func(ctx *Ctx) (models.Message, error)

	Ctx struct {
		bot      *BotAPI
		Endpoint string
		Message  *models.Message
		Params   map[string]interface{}
		context.Context
		HandlerIndex uint
		UserState    UserState
	}

	State struct {
		Endpoint string
		Message  *models.Message
		Params   map[string]interface{}
	}

	UserState struct {
		Stack []State
		Next  *State
	}
)

func (ctx *Ctx) Bot() *BotAPI {
	return ctx.bot
}

func (ctx *Ctx) WithParam(key string, val interface{}) *Ctx {
	ctx.Params[key] = val
	return ctx
}

func (ctx *Ctx) GetParam(key string) interface{} {
	return ctx.Params[key]
}

func (ctx *Ctx) WithContext(newCtx context.Context) {
	ctx.Context = newCtx
}

func (ctx *Ctx) GetContext() context.Context {
	return ctx.Context
}

func (ctx *Ctx) ResetUserStack() {
	ctx.bot.userStats[ctx.Message.From.Id] = UserState{
		Stack: make([]State, 0),
		Next:  nil,
	}
}

func (ctx *Ctx) SetNextStat(state State) {
	userState := ctx.bot.userStats[ctx.Message.From.Id]
	userState.Next = &state
	ctx.bot.userStats[ctx.Message.From.Id] = userState
}
