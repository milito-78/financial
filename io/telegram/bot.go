package telegram

import (
	"context"
	"errors"
	"financial/application"
	"financial/domain"
	"financial/infrastructure/cache"
	telBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"net/url"
	"strconv"
)

type TlgBot struct {
	userService application.IUserService
	bot         *telBot.BotAPI
	router      *Router
	state       StateManager
	ctx         context.Context
}

func NewTelegramBot(ctx context.Context, cache cache.ICache, token string, userService application.IUserService) *TlgBot {
	defer log.Info("Bot created successfully")
	bot, err := telBot.NewBotAPI(token)
	if err != nil {
		log.Fatalf("error during new bot : %s \n", err)
	}
	log.Info("Creating new bot...")
	return &TlgBot{bot: bot, router: NewRouter(), state: StateManager{cache: cache}, ctx: ctx, userService: userService}
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
		message := update.Message
		var err error

		user, err := tlg.getCurrentUser(update.SentFrom())
		if err != nil {
			handled := handleErrorView(err)
			tlg.bot.Send(handled.Render(update))
			continue
		}

		var lastState *State
		var route *Route
		var params []string
		var queries *url.Values
		var path string

		switch {
		case message != nil && tlg.onlyTextMessage(message):
			log.Info("Receive text message")
			lastState, err = tlg.handleLastStateForUser(strconv.FormatInt(message.From.ID, 10))
			path = message.Text
		case update.CallbackQuery != nil:
			log.Info("Receive callback message")
			lastState, err = tlg.handleLastStateForUser(strconv.FormatInt(update.SentFrom().ID, 10))
			path = update.CallbackQuery.Data
		default:
			log.Info("Unsupported message")
			handled := handleErrorView(&RouteNotFoundError{})
			tlg.bot.Send(handled.Render(update))
			continue
		}
		route, queries, params, err = tlg.router.MatchRoute(path)
		if err != nil && (lastState == nil || lastState.Data.(string) == "") {
			handled := handleErrorView(err)
			tlg.bot.Send(handled.Render(update))
			continue
		} else if route == nil {
			route, queries, params, err = tlg.router.MatchRoute(lastState.Data.(string))
			if err != nil {
				handled := handleErrorView(DataNotFound{})
				tlg.bot.Send(handled.Render(update))
				continue
			}
		}
		ctx := &RequestContext{
			RouteParams:       params,
			QueryParams:       *queries,
			Received:          update,
			Route:             route,
			Message:           path,
			stateSaverChannel: saverCh,
			lastState:         lastState,
			user:              user,
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

func (tlg *TlgBot) onlyTextMessage(message *telBot.Message) bool {
	return message.Text != "" && message.Photo == nil && message.Document == nil && message.Video == nil
}

func (tlg *TlgBot) registerRouter() {
	registerRoutes(tlg.Router())
}

func (tlg *TlgBot) renderView(ctx *RequestContext, view Renderable) {
	send, err := tlg.bot.Send(view.Render(ctx.Received))
	log.Printf("Send message : %s %v", err, send)
}

func (tlg *TlgBot) handleLastStateForUser(lastStateUser string) (*State, error) {
	found, err := tlg.state.Get(tlg.ctx, lastStateUser)
	if err != nil {
		return nil, err
	}
	return found, nil
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
				if res.Data == nil {
					tlg.state.Delete(tlg.ctx, res.Uid)
				} else {
					tlg.state.Set(tlg.ctx, res)
				}
			}
		}
	}()
	return ch
}

func (tlg *TlgBot) getCurrentUser(telUser *telBot.User) (*domain.User, error) {
	uid := strconv.FormatInt(telUser.ID, 10)
	user, err := tlg.userService.GetUserByUuid(uid)
	if err != nil {
		if err != nil {
			err := tlg.userService.AddUser(domain.NewUser(telUser.UserName, telUser.FirstName, telUser.LastName, uid))
			if err != nil {
				log.Printf("Error during create new user : %s \n", err)
				return nil, err
			}
		}
	}

	return user, nil
}

func handleErrorView(err error) Renderable {
	if errors.As(err, &DataNotFound{}) {
		return NewNotFoundView(err.Error())
	}
	if errors.As(err, &RouteNotFoundError{}) {
		return NewNotFoundView(err.Error())
	}
	if errors.As(err, &AccessError{}) {
		return NewAccessErrorView(err.Error())
	}

	//unknown error
	return NewUnknownErrorView()
}
