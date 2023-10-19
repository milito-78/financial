package telegram

import (
	"errors"
	"financial/application"
	"financial/domain"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
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
	user := ctx.GetUser()
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
	groupId, err := strconv.ParseUint(ctx.RouteParams[0], 10, 64)
	if err != nil {
		return nil, DataNotFound{}
	}

	result, err := g.groupService.GetGroup(groupId)
	if err != nil {
		return nil, DataNotFound{}
	}

	return NewGroupShowView(result, ctx.GetUser()), nil
}

func (g *GroupActionCommand) Create(ctx *RequestContext) (Renderable, error) {
	user := ctx.GetUser()
	ctx.SetState(State{user.Uuid, GroupsStoreName})
	return NewGroupCreateView(), nil
}

func (g *GroupActionCommand) Store(ctx *RequestContext) (Renderable, error) {
	user := ctx.GetUser()
	gUuid := uuid.New()

	var group = &domain.Group{
		InviteLink: gUuid.String(),
		CreatorId:  user.ID,
		Name:       ctx.Message,
	}

	current := ctx.GetState()
	switch current.Data.(string) {
	case GroupsStoreName:
		err := g.groupService.Store(group)
		if err != nil {
			log.Printf("Error during create group : %s", err)
			return nil, UnknownError{}
		}
	default:
		return nil, RouteNotFoundError{}
	}
	ctx.SetState(State{user.Uuid, nil})
	return NewGroupStoreView(group), nil
}

func (g *GroupActionCommand) Edit(ctx *RequestContext) (Renderable, error) {
	user := ctx.GetUser()
	groupId, _ := strconv.ParseUint(ctx.RouteParams[0], 10, 16)
	group, err := g.groupService.GetGroup(groupId)
	if err != nil {
		return nil, err
	}
	var text string
	action := ctx.QueryParams.Get("action")
	updateRoute := strings.Replace(GroupsUpdate, ":id", strconv.FormatUint(group.ID, 10), 1) + "?action=" + action
	switch action {
	case RouteParamEditGroupName:
		text = fmt.Sprintf(`
You are changing %s group name.
Please enter your new name for exisits group : 
`, group.Name)
		ctx.SetState(State{Uid: user.Uuid, Data: updateRoute})
	case RouteParamEditGroupLink:
		text = fmt.Sprintf(`
You are changing %s group invite link.
Do you want to change your group invite link? 
`, group.Name)
		ctx.SetState(State{Uid: user.Uuid, Data: updateRoute})
	default:
		return nil, errors.New("action is invalid")
	}

	return NewGroupEditView(text, group, action), nil
}

func (g *GroupActionCommand) Update(ctx *RequestContext) (Renderable, error) {
	user := ctx.GetUser()
	defer ctx.SetState(State{user.Uuid, nil})

	groupId, _ := strconv.ParseUint(ctx.RouteParams[0], 10, 16)
	group, err := g.groupService.GetGroup(groupId)
	if err != nil {
		return nil, err
	}
	var text string
	action := ctx.QueryParams.Get("action")
	switch action {
	case RouteParamEditGroupName:
		group.Name = ctx.Message
		text = "Group name updated successfully."
	case RouteParamEditGroupLink:
		confirmed := ctx.QueryParams.Get("confirm")
		if confirmed != "1" {
			return nil, RouteNotFoundError{}
		}

		gUuid := uuid.New()
		group.InviteLink = gUuid.String()
		text = "Group invite link changed successfully."
	default:
		return nil, errors.New("action is invalid")
	}

	err = g.groupService.Update(group)
	if err != nil {
		return nil, err
	}

	return NewGroupUpdateView(text, group), nil
}

func (g *GroupActionCommand) Delete(ctx *RequestContext) (Renderable, error) {
	user := ctx.GetUser()
	groupId, _ := strconv.ParseUint(ctx.RouteParams[0], 10, 16)
	group, err := g.groupService.GetGroup(groupId)
	if err != nil {
		return nil, err
	}
	if group.CreatorId != user.ID {
		return nil, AccessError{ActionName: "Delete group"}
	}

	var text string
	var confirm = false
	if ctx.QueryParams.Get("confirm") == "1" {
		if res := g.groupService.DeleteGroup(group.ID); res == true {
			text = "group deleted successfully"
		} else {
			return nil, UnknownError{}
		}
	} else {
		text = fmt.Sprintf("Do you want to delete your group called %s? If you delete it you cannot get back.", group.Name)
		confirm = true
	}

	return NewGroupDeleteView(text, group, confirm), nil
}

func (g *GroupActionCommand) Leave(ctx *RequestContext) (Renderable, error) {
	user := ctx.GetUser()
	groupId, _ := strconv.ParseUint(ctx.RouteParams[0], 10, 16)
	group, err := g.groupService.GetGroup(groupId)
	if err != nil {
		return nil, err
	}
	if group.CreatorId == user.ID {
		return nil, AccessError{ActionName: "Leave group"}
	}

	var text string
	var confirm = false
	if ctx.QueryParams.Get("confirm") == "1" {
		if res := g.groupService.LeaveGroup(group.ID, user.ID); res == true {
			text = "You leaved group successfully"
		} else {
			return nil, UnknownError{}
		}
	} else {
		text = fmt.Sprintf("Do you want to leave your group called %s?", group.Name)
		confirm = true
	}

	return NewGroupDeleteView(text, group, confirm), nil
}
