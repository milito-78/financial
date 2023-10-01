package telegram

import (
	"financial/domain"
	"financial/infrastructure/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

type ButtonComponent struct {
	Title string
}

type InlineButtonComponent struct {
	Title    string
	CallBack string
}

type Renderable interface {
	Render(received tgbotapi.Update) tgbotapi.MessageConfig
}

type NotFoundView struct {
	Text string
}

func NewNotFoundView(text string) *NotFoundView {
	if text == "" {
		text = "Command is not found."
	}
	return &NotFoundView{Text: text}
}

func (s NotFoundView) Render(received tgbotapi.Update) tgbotapi.MessageConfig {
	message := tgbotapi.NewMessage(received.SentFrom().ID, s.Text)
	return message
}

type UnknownErrorView struct {
}

func NewUnknownErrorView() *UnknownErrorView {
	return &UnknownErrorView{}
}

func (s UnknownErrorView) Render(received tgbotapi.Update) tgbotapi.MessageConfig {
	message := tgbotapi.NewMessage(received.SentFrom().ID, "There is an error during handle your request. Please be patient and try again later.")
	return message
}

type StartView struct {
	Text    string
	Buttons []ButtonComponent
}

func NewStartView(name string, groupStatus, paymentStatus, myPaymentsStatus bool, startCmd bool) *StartView {
	text := "How can I help you dear " + name
	if startCmd {
		text = "Hello dear " + name + ". Welcome to payment handler bot.\n" +
			"This bot helps you to share payments and invoices between your friends."
	}

	buttons := []ButtonComponent{
		{Title: GroupsCmd},
		{Title: "Invoices menu"},
		{Title: "My Invoices"},
	}
	return &StartView{Text: text, Buttons: buttons}
}

func (s StartView) Render(received tgbotapi.Update) tgbotapi.MessageConfig {
	message := tgbotapi.NewMessage(received.Message.Chat.ID, s.Text)
	var columns = buttonMakerMatrix(s.Buttons, 2)
	replyKeyboard := tgbotapi.NewReplyKeyboard(columns...)
	message.ReplyMarkup = replyKeyboard
	return message
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
	Paginate db.Paginate[domain.Group]
}

func NewGroupListView(paginate db.Paginate[domain.Group]) *GroupListView {
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
		text = "Here is your groups "
		for _, result := range g.Paginate.Results {
			inlineKeyboards = append(inlineKeyboards, InlineButtonComponent{
				Title:    strconv.FormatUint(result.ID, 10) + result.Name,
				CallBack: strings.Replace(GroupsShow, ":id", strconv.FormatUint(result.ID, 10), 1),
			})
		}
	}

	var inlineColumn = inlineButtonMakerMatrix(inlineKeyboards, 3)
	inlineColumn = append(inlineColumn, []tgbotapi.InlineKeyboardButton{
		//tgbotapi.NewInlineKeyboardButtonData("Next Page", "1"),
		//tgbotapi.NewInlineKeyboardButtonData("Prev Page", "2"),
	})
	message := tgbotapi.NewMessage(received.Message.Chat.ID, text)
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
	message := tgbotapi.NewMessage(received.Message.Chat.ID, text)
	return message
}

func buttonMakerMatrix(buttons []ButtonComponent, colCount uint) [][]tgbotapi.KeyboardButton {
	length := len(buttons)
	if length == 0 {
		return nil
	}

	j := int(colCount)
	i := length / j
	if i == 0 {
		i = 0
	}

	if length%j != 0 {
		i += 1
	}

	tlgButtons := make([][]tgbotapi.KeyboardButton, i)
	for k := range tlgButtons {
		tlgButtons[k] = make([]tgbotapi.KeyboardButton, j)
	}

	row, col := 0, 0
	for _, btn := range buttons {
		tlgButtons[row][col] = tgbotapi.NewKeyboardButton(btn.Title)
		if col+1 < j {
			col += 1
			continue
		} else {
			col = 0
			row += 1
		}
	}

	return tlgButtons
}

func inlineButtonMakerMatrix(buttons []InlineButtonComponent, colCount uint) [][]tgbotapi.InlineKeyboardButton {
	length := len(buttons)
	if length == 0 {
		return nil
	}

	j := int(colCount)
	i := length / j

	if length%j != 0 {
		i += 1
	}

	tlgButtons := make([][]tgbotapi.InlineKeyboardButton, i)
	for k := range tlgButtons {
		tlgButtons[k] = make([]tgbotapi.InlineKeyboardButton, j)
	}

	row, col := 0, 0
	for _, btn := range buttons {
		tlgButtons[row][col] = tgbotapi.NewInlineKeyboardButtonData(btn.Title, btn.CallBack)
		if col+1 < j {
			col += 1
			continue
		} else {
			col = 0
			row += 1
		}
	}

	newOne := make([][]tgbotapi.InlineKeyboardButton, i)
	for x, buttons := range tlgButtons {
		for _, button := range buttons {
			if button.Text != "" {
				newOne[x] = append(newOne[x], button)
			}
		}
	}

	return newOne
}
