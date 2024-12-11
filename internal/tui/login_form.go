package tui

import (
	"github.com/moonicy/goph-keeper-yandex/internal/events"
	"github.com/moonicy/goph-keeper-yandex/internal/grpc_client"
	"github.com/moonicy/goph-keeper-yandex/internal/subscribtion"
	"github.com/rivo/tview"
)

type LoginForm struct {
	sub     *subscribtion.Subscription
	client  *grpc_client.Client
	pages   *tview.Pages
	message *Message
}

func NewLoginForm(sub *subscribtion.Subscription, client *grpc_client.Client, pages *tview.Pages, message *Message) *LoginForm {
	loginForm := &LoginForm{
		sub:     sub,
		client:  client,
		pages:   pages,
		message: message,
	}
	sub.SubscribeEvent(events.ShowLoginForm, loginForm.GetNotifier())
	return loginForm
}

// Show Функция для отображения формы авторизации
func (lf *LoginForm) Show() {
	var loginForm *tview.Form

	loginForm = tview.NewForm().
		AddInputField("Логин", "", 20, nil, nil).
		AddPasswordField("Пароль", "", 20, '*', nil).
		AddButton("Войти", func() {
			login := loginForm.GetFormItemByLabel("Логин").(*tview.InputField).GetText()
			password := loginForm.GetFormItemByLabel("Пароль").(*tview.InputField).GetText()

			err := lf.client.Login(login, password)
			if err != nil {
				lf.message.ShowError(err)
				return
			}

			// Переходим в главное меню
			lf.sub.NotifyEvent(events.ShowMainMenu, nil)
		}).
		AddButton("Зарегистрироваться", func() {
			lf.sub.NotifyEvent(events.ShowRegisterForm, nil)
		})
	loginForm.SetTitle("Авторизация").SetBorder(true)

	lf.pages.AddAndSwitchToPage("login", loginForm, true)
}

func (lf *LoginForm) GetNotifier() subscribtion.Notifier {
	return func(_ any) {
		lf.Show()
	}
}
