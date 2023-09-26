package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	message := tgbotapi.NewMessage(received.Message.Chat.ID, s.Text)
	return message
}

type UnknownErrorView struct {
}

func NewUnknownErrorView() *UnknownErrorView {
	return &UnknownErrorView{}
}

func (s UnknownErrorView) Render(received tgbotapi.Update) tgbotapi.MessageConfig {
	message := tgbotapi.NewMessage(received.Message.Chat.ID, "There is an error during handle your request. Please be patient and try again later.")
	return message
}

type StartView struct {
	Text    string
	Buttons []ButtonComponent
}

func NewStartView(name string, groupStatus, paymentStatus, myPaymentsStatus bool) *StartView {
	text := "Hello dear " + name + ". Welcome to payment handler bot.\n" +
		"This bot helps you to share payments and invoices between your friends."
	buttons := []ButtonComponent{
		{Title: "Groups menu"},
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
