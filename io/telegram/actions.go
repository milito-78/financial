package telegram

import (
	"financial/application"
	"financial/domain"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type Handler interface {
	Handle(ctx *RequestContext) (Renderable, error)
}

type StartCommand struct {
	userService application.IUserService
}

func NewStartCommand(userService application.IUserService) *StartCommand {
	return &StartCommand{userService: userService}
}

func (s *StartCommand) Handle(ctx *RequestContext) (Renderable, error) {
	telUser := ctx.Received.SentFrom()
	uuid := strconv.FormatInt(telUser.ID, 10)
	_, err := s.userService.GetUserByUuid(uuid)
	if err != nil {
		err := s.userService.AddUser(domain.NewUser(telUser.UserName, telUser.FirstName, telUser.LastName, uuid))
		if err != nil {
			log.Printf("Error during create new user : %s \n", err)
		}
	}

	return NewStartView(telUser.UserName, true, true, true), nil
}
