package telegram

import (
	"context"
	"errors"
	"financial/infrastructure/cache"
	telBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

type TlgBot struct {
	bot    *telBot.BotAPI
	router *Router
	state  StateManager
	ctx    context.Context
}

func NewTelegramBot(ctx context.Context, cache cache.ICache, token string) *TlgBot {
	defer log.Info("Bot created successfully")
	bot, err := telBot.NewBotAPI(token)
	if err != nil {
		log.Fatalf("error during new bot : %s \n", err)
	}
	log.Info("Creating new bot...")
	return &TlgBot{bot: bot, router: NewRouter(), state: StateManager{cache: cache}, ctx: ctx}
}

func (tlg *TlgBot) StartBot() {
	log.Info("Starting bot")
	saverCh := tlg.startStateChannel()
	log.Info("Registering routes...")
	tlg.registerRouter()
	log.Info("Registering finished")

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
			RouteParams:       params,
			Received:          update,
			Route:             route,
			Message:           message,
			StateSaverChannel: saverCh,
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

func (tlg *TlgBot) registerRouter() {
	registerRoutes(tlg.Router())
}

func (tlg *TlgBot) renderView(ctx *RequestContext, view Renderable) {
	tlg.bot.Send(view.Render(ctx.Received))
}

func (tlg *TlgBot) handleLastStateForUser(lastStateUser string) (string, error) {
	found, err := tlg.state.Get(tlg.ctx, lastStateUser)
	if err != nil {
		return "", err
	}
	return found.Data.(string), nil
}

func (tlg *TlgBot) startStateChannel() chan<- State {
	ch := make(chan State)
	go func() {
		for {
			select {
			case <-tlg.ctx.Done():
				log.Warn("Shutting down the redis client")
				return
			case res := <-ch:
				tlg.state.Set(tlg.ctx, res)
			}
		}
	}()
	return ch
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
