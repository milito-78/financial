package telegram

import (
	"financial/config"
	"financial/domain"
	"financial/infrastructure/db"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

type Renderable interface {
	Render(received tgbotapi.Update) tgbotapi.MessageConfig
}

type GroupMenuView struct {
	Text    string
	Buttons []ButtonComponent
}

func NewGroupMenuView() *GroupMenuView {
	text := "This is group menu. Please select your request from keyboard:"
	buttons := []ButtonComponent{
		{Title: GroupsList},
		{Title: "Create New One"},
		{Title: "Back to menu"},
	}
	return &GroupMenuView{Text: text, Buttons: buttons}
}

func (s GroupMenuView) Render(received tgbotapi.Update) tgbotapi.MessageConfig {
	message := tgbotapi.NewMessage(received.Message.Chat.ID, s.Text)
	var columns = buttonMakerMatrix(s.Buttons, 2)
	replyKeyboard := tgbotapi.NewReplyKeyboard(columns...)
	message.ReplyMarkup = replyKeyboard
	return message
}

type GroupListView struct {
	Paginate *db.Paginate[domain.Group]
}

func NewGroupListView(paginate *db.Paginate[domain.Group]) *GroupListView {
	return &GroupListView{Paginate: paginate}
}

func (g GroupListView) Render(received tgbotapi.Update) tgbotapi.MessageConfig {
	var text string
	var inlineKeyboards []InlineButtonComponent
	if len(g.Paginate.Results) == 0 {
		text = "You dont have any group. Please start with creating one!"
		inlineKeyboards = append(inlineKeyboards, InlineButtonComponent{
			Title:    "Create New One!",
			CallBack: GroupsCreate,
		})
	} else {
		text = "Groups that you created: "
		for _, result := range g.Paginate.Results {
			inlineKeyboards = append(inlineKeyboards, InlineButtonComponent{
				Title:    fmt.Sprintf("%d. %s", result.ID, result.Name),
				CallBack: strings.Replace(GroupsShow, ":id", strconv.FormatUint(result.ID, 10), 1),
			})
		}
	}

	var inlineColumn = inlineButtonMakerMatrix(inlineKeyboards, 3)
	var temp []tgbotapi.InlineKeyboardButton
	if g.Paginate.Page != 1 {
		temp = append(temp, tgbotapi.NewInlineKeyboardButtonData("Prev Page", GroupsList+"?page="+strconv.Itoa(int(g.Paginate.Page-1))))
	}
	if g.Paginate.NextPage {
		temp = append(temp, tgbotapi.NewInlineKeyboardButtonData("Next Page", GroupsList+"?page="+strconv.Itoa(int(g.Paginate.Page+1))))
	}

	if len(temp) != 0 {
		inlineColumn = append(inlineColumn, temp)
	}

	message := tgbotapi.NewMessage(received.SentFrom().ID, text)
	message.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(inlineColumn...)

	return message
}

type GroupCreateView struct {
}

func NewGroupCreateView() *GroupCreateView {
	return &GroupCreateView{}
}

func (g GroupCreateView) Render(received tgbotapi.Update) tgbotapi.MessageConfig {
	var text string
	text = "Please tell me what is your group name? :"
	message := tgbotapi.NewMessage(received.SentFrom().ID, text)
	return message
}

type GroupStoreView struct {
	group *domain.Group
}

func NewGroupStoreView(group *domain.Group) *GroupStoreView {
	return &GroupStoreView{group: group}
}

func (g GroupStoreView) Render(received tgbotapi.Update) tgbotapi.MessageConfig {
	text := fmt.Sprintf(`
Group created successfully, Here group details:
<b>Group Name :</b> %s
<a href="https://t.me/%s?start=%s">Invite link</a>
`, g.group.Name, config.Default.(*config.App).BotId, g.group.InviteLink)

	message := tgbotapi.NewMessage(received.SentFrom().ID, text)
	message.ParseMode = tgbotapi.ModeHTML
	return message
}

type GroupShowView struct {
	group *domain.Group
	user  *domain.User
}

func NewGroupShowView(group *domain.Group, user *domain.User) *GroupShowView {
	return &GroupShowView{group: group, user: user}
}

func (g GroupShowView) Render(received tgbotapi.Update) tgbotapi.MessageConfig {
	text := fmt.Sprintf(`
Here group details:
<b>Group Name :</b> %s
<b>Owner Username :</b> @%s
`, g.group.Name, g.group.Creator.Username)
	if g.user.ID == g.group.CreatorId {
		text = text + fmt.Sprintf("<a href=\"https://t.me/%s?start=%s\">Invite link</a>", config.Default.(*config.App).BotId, g.group.InviteLink)
	} else {
		text = text + "You were invite to this group."
	}

	message := tgbotapi.NewMessage(received.SentFrom().ID, text)

	if g.user.ID == g.group.CreatorId {
		editRoute := strings.Replace(GroupsEdit, ":id", strconv.FormatUint(g.group.ID, 10), 1)
		message.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData("Edit Name", editRoute+"?action="+RouteParamEditGroupName),
				tgbotapi.NewInlineKeyboardButtonData("Change Invite Link", editRoute+"?action="+RouteParamEditGroupLink),
			},
			[]tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData("Delete", strings.Replace(GroupsDelete, ":id", strconv.FormatUint(g.group.ID, 10), 1)),
			},
			[]tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData("Members", strings.Replace(GroupsMembers, ":id", strconv.FormatUint(g.group.ID, 10), 1)),
			},
		)
	} else {
		message.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData("Leave Group", strings.Replace(GroupsLeave, ":id", strconv.FormatUint(g.group.ID, 10), 1)),
			},
		)
	}

	message.ParseMode = tgbotapi.ModeHTML
	return message
}

type GroupEditView struct {
	group  *domain.Group
	action string
	text   string
}

func NewGroupEditView(text string, group *domain.Group, action string) *GroupEditView {
	return &GroupEditView{group: group, action: action, text: text}
}

func (g GroupEditView) Render(received tgbotapi.Update) tgbotapi.MessageConfig {
	message := tgbotapi.NewMessage(received.SentFrom().ID, g.text)
	message.ParseMode = tgbotapi.ModeHTML
	if g.action == RouteParamEditGroupLink {
		editRoute := strings.Replace(GroupsUpdate, ":id", strconv.FormatUint(g.group.ID, 10), 1)
		message.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData("Yes, change it", editRoute+"?action="+RouteParamEditGroupLink+"&confirm=1"),
				tgbotapi.NewInlineKeyboardButtonData("No, Came back to list", GroupsList),
			},
		)
	}

	return message
}

type GroupUpdateView struct {
	group *domain.Group
	text  string
}

func NewGroupUpdateView(text string, group *domain.Group) *GroupUpdateView {
	return &GroupUpdateView{group: group, text: text}
}

func (g GroupUpdateView) Render(received tgbotapi.Update) tgbotapi.MessageConfig {
	text := g.text
	text = fmt.Sprintf(`
Here group details:
<b>Group Name :</b> %s
<b>Owner Username :</b> @%s
<a href="https://t.me/%s?start=%s">Invite link</a>
`, g.group.Name, g.group.Creator.Username, config.Default.(*config.App).BotId, g.group.InviteLink)
	message := tgbotapi.NewMessage(received.SentFrom().ID, text)
	message.ParseMode = tgbotapi.ModeHTML

	return message
}

type GroupDeleteView struct {
	group   *domain.Group
	text    string
	confirm bool
}

func NewGroupDeleteView(text string, group *domain.Group, confirm bool) *GroupDeleteView {
	return &GroupDeleteView{group: group, text: text, confirm: confirm}
}

func (g GroupDeleteView) Render(received tgbotapi.Update) tgbotapi.MessageConfig {
	text := g.text
	message := tgbotapi.NewMessage(received.SentFrom().ID, text)

	if g.confirm {
		message.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData("Yes, delete it", strings.Replace(GroupsDelete, ":id", strconv.FormatUint(g.group.ID, 10), 1)+"?confirm=1"),
				tgbotapi.NewInlineKeyboardButtonData("No, Came back to list", GroupsList),
			},
		)
	}

	return message
}

type GroupLeaveView struct {
	group   *domain.Group
	text    string
	confirm bool
}

func NewGroupLeaveView(text string, group *domain.Group, confirm bool) *GroupLeaveView {
	return &GroupLeaveView{group: group, text: text, confirm: confirm}
}

func (g GroupLeaveView) Render(received tgbotapi.Update) tgbotapi.MessageConfig {
	text := g.text
	message := tgbotapi.NewMessage(received.SentFrom().ID, text)

	if g.confirm {
		message.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData("Yes, I want to leave", strings.Replace(GroupsLeave, ":id", strconv.FormatUint(g.group.ID, 10), 1)+"?confirm=1"),
				tgbotapi.NewInlineKeyboardButtonData("No, Came back to list", GroupsList),
			},
		)
	}

	return message
}
