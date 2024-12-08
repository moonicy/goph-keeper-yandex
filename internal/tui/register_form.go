package tui

import (
	"fmt"
	"github.com/moonicy/goph-keeper-yandex/internal/events"
	"github.com/moonicy/goph-keeper-yandex/internal/grpc_client"
	"github.com/moonicy/goph-keeper-yandex/internal/subscribtion"
	"github.com/rivo/tview"
)

type RegisterForm struct {
	sub     *subscribtion.Subscription
	client  *grpc_client.Client
	pages   *tview.Pages
	message *Message
}

func NewRegisterForm(sub *subscribtion.Subscription, client *grpc_client.Client, pages *tview.Pages, message *Message) *RegisterForm {
	registerForm := &RegisterForm{
		sub:     sub,
		client:  client,
		pages:   pages,
		message: message,
	}
	sub.SubscribeEvent(events.ShowRegisterForm, registerForm.GetNotifier())
	return registerForm
}

// Show Функция для отображения формы регистрации
func (rf *RegisterForm) Show() {
	var registerForm *tview.Form

	registerForm = tview.NewForm().
		AddInputField("Логин", "", 20, nil, nil).
		AddPasswordField("Пароль", "", 20, '*', nil).
		AddButton("Зарегистрироваться", func() {
			login := registerForm.GetFormItemByLabel("Логин").(*tview.InputField).GetText()
			password := registerForm.GetFormItemByLabel("Пароль").(*tview.InputField).GetText()

			userID, err := rf.client.Register(login, password)
			if err != nil {
				rf.message.ShowError(err)
				return
			}

			rf.message.ShowMessage(fmt.Sprintf("Успешно зарегистрированы. Ваш ID: %d", userID))

			// Возвращаемся к форме авторизации
			rf.sub.NotifyEvent(events.ShowLoginForm, nil)
		}).
		AddButton("Назад", func() {
			rf.sub.NotifyEvent(events.ShowLoginForm, nil)
		})
	registerForm.SetTitle("Регистрация").SetBorder(true)

	rf.pages.AddAndSwitchToPage("register", registerForm, true)
}

func (rf *RegisterForm) GetNotifier() subscribtion.Notifier {
	return func(_ any) {
		rf.Show()
	}
}
