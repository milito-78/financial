package telegram

import (
	"financial/application"
	"financial/domain"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type StartCommand struct {
	userService application.IUserService
}

func NewStartCommand(userService application.IUserService) *StartCommand {
	return &StartCommand{userService: userService}
}

func (s *StartCommand) Handle(ctx *RequestContext) (Renderable, error) {
	telUser := ctx.Received.SentFrom()
	isStart := false
	if ctx.Route.Path != BackToMenu {
		isStart = true
		uuid := strconv.FormatInt(telUser.ID, 10)
		_, err := s.userService.GetUserByUuid(uuid)
		if err != nil {
			err := s.userService.AddUser(domain.NewUser(telUser.UserName, telUser.FirstName, telUser.LastName, uuid))
			if err != nil {
				log.Printf("Error during create new user : %s \n", err)
			}
		}
	}

	return NewStartView(telUser.UserName, true, true, true, isStart), nil
}

type GroupActionCommand struct {
	userService  application.IUserService
	groupService application.IGroupService
}

func NewGroupActionCommand(userService application.IUserService, groupService application.IGroupService) *GroupActionCommand {
	return &GroupActionCommand{userService: userService, groupService: groupService}
}

func (g *GroupActionCommand) MenuHandle(ctx *RequestContext) (Renderable, error) {
	return NewGroupMenuView(), nil
}

func (g *GroupActionCommand) List(ctx *RequestContext) (Renderable, error) {
	telUser := ctx.Received.SentFrom()
	uuid := strconv.FormatInt(telUser.ID, 10)
	user, err := g.userService.GetUserByUuid(uuid)
	if err != nil {
		return nil, err
	}

	result := g.groupService.UserGroupPaginate(user.ID, 1)
	return NewGroupListView(result), nil
}

func (g *GroupActionCommand) Show(ctx *RequestContext) (Renderable, error) {
	telUser := ctx.Received.SentFrom()
	uuid := strconv.FormatInt(telUser.ID, 10)
	user, err := g.userService.GetUserByUuid(uuid)
	if err != nil {
		return nil, err
	}

	result := g.groupService.UserGroupPaginate(user.ID, 1)
	return NewGroupListView(result), nil
}

func (g *GroupActionCommand) Create(ctx *RequestContext) (Renderable, error) {
	telUser := ctx.Received.SentFrom()
	uuid := strconv.FormatInt(telUser.ID, 10)
	_, err := g.userService.GetUserByUuid(uuid)
	if err != nil {
		return nil, err
	}

	return NewGroupCreateView(), nil
}
