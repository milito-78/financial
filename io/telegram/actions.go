package telegram

import (
	"financial/application"
	"financial/domain"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"strconv"
)

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
	uid := strconv.FormatInt(telUser.ID, 10)
	user, err := g.userService.GetUserByUuid(uid)
	if err != nil {
		return nil, err
	}
	page := uint(1)
	if tmp := ctx.QueryParams.Get("page"); tmp != "" {
		if t, err := strconv.Atoi(tmp); err == nil {
			page = uint(t)
		}
	}
	result := g.groupService.UserGroupPaginate(user.ID, page)
	return NewGroupListView(result), nil
}

func (g *GroupActionCommand) Show(ctx *RequestContext) (Renderable, error) {
	telUser := ctx.Received.SentFrom()
	uid := strconv.FormatInt(telUser.ID, 10)
	user, err := g.userService.GetUserByUuid(uid)
	if err != nil {
		return nil, err
	}

	result, _ := g.groupService.UserGetGroup(user.ID, 1)
	return NewGroupShowView(result), nil
}

func (g *GroupActionCommand) Create(ctx *RequestContext) (Renderable, error) {
	telUser := ctx.Received.SentFrom()
	uid := strconv.FormatInt(telUser.ID, 10)
	_, err := g.userService.GetUserByUuid(uid)
	if err != nil {
		return nil, err
	}
	ctx.SetState(State{uid, GroupsStoreName})
	return NewGroupCreateView(), nil
}

func (g *GroupActionCommand) Store(ctx *RequestContext) (Renderable, error) {
	telUser := ctx.Received.SentFrom()
	uid := strconv.FormatInt(telUser.ID, 10)
	user, err := g.userService.GetUserByUuid(uid)
	if err != nil {
		return nil, err
	}

	gUuid := uuid.New()

	var group = &domain.Group{
		InviteLink: gUuid.String(),
		CreatorId:  user.ID,
		Name:       ctx.Message,
	}

	current := ctx.GetState()
	switch current.Data.(string) {
	case GroupsStoreName:
		err = g.groupService.Store(group)
		if err != nil {
			log.Printf("Error during create group : %s", err)
			return nil, UnknownError{}
		}
	default:
		return nil, RouteNotFoundError{}
	}
	ctx.SetState(State{uid, nil})
	return NewGroupStoreView(group), nil
}
