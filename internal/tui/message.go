package tui

import "github.com/rivo/tview"

type Message struct {
	pages *tview.Pages
}

func NewMessage(pages *tview.Pages) *Message {
	return &Message{
		pages: pages,
	}
}

// ShowError Функция для отображения ошибок
func (m *Message) ShowError(err error) {
	modal := tview.NewModal().
		SetText(err.Error()).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			m.pages.RemovePage("modal")
		})
	m.pages.AddPage("modal", modal, false, true)
}

// ShowMessage Функция для отображения сообщений
func (m *Message) ShowMessage(message string) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			m.pages.RemovePage("modal")
		})
	m.pages.AddPage("modal", modal, false, true)
}
