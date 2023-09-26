package telegram

import (
	"errors"
	telBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

type TlgBot struct {
	bot    *telBot.BotAPI
	router *Router
}

func NewTelegramBot(token string) *TlgBot {
	defer log.Info("Bot created successfully")
	bot, err := telBot.NewBotAPI(token)
	if err != nil {
		log.Fatalf("error during new bot : %s \n", err)
	}
	log.Info("Creating new bot...")
	return &TlgBot{bot: bot, router: NewRouter()}
}

func (tlg *TlgBot) StartBot() {
	log.Info("Starting bot")
	updateConfig := telBot.NewUpdate(0)

	updates := tlg.bot.GetUpdatesChan(updateConfig)
	log.Info("Waiting for new update")
	for update := range updates {
		route, message, params, err := tlg.router.MatchRoute(update)
		if err != nil {
			handled := handleErrorView(err)
			tlg.bot.Send(handled.Render(update))
			continue
		}

		ctx := &RequestContext{
			RouteParams: params,
			Received:    update,
			Route:       route,
			Message:     message,
		}
		renderable, err := route.Handler(ctx)
		if err != nil {
			handled := handleErrorView(err)
			tlg.bot.Send(handled.Render(update))
			continue
		}
		tlg.renderView(ctx, renderable)
	}
}

func (tlg *TlgBot) StopBot() {
	tlg.bot.StopReceivingUpdates()
}

func (tlg *TlgBot) Router() *Router {
	return tlg.router
}

func (tlg *TlgBot) renderView(ctx *RequestContext, view Renderable) {
	tlg.bot.Send(view.Render(ctx.Received))
}

func handleErrorView(err error) Renderable {
	if errors.Is(err, DataNotFound{}) {
		return NewNotFoundView(err.Error())
	}
	if errors.Is(err, RouteNotFoundError{}) {
		return NewNotFoundView("")
	}

	//unknown error
	return NewUnknownErrorView()
}
