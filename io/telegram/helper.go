package telegram

import (
	"financial/application"
	"financial/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	//TODO handle invite to group and make them friend
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

// Buttons
type ButtonComponent struct {
	Title string
}

type InlineButtonComponent struct {
	Title    string
	CallBack string
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

// Views
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
