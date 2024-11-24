package tui

import (
	"github.com/moonicy/goph-keeper-yandex/internal/events"
	"github.com/moonicy/goph-keeper-yandex/internal/grpc_client"
	"github.com/moonicy/goph-keeper-yandex/internal/subscribtion"
	"github.com/rivo/tview"
)

type AddDataMenu struct {
	sub     *subscribtion.Subscription
	client  *grpc_client.Client
	pages   *tview.Pages
	message *Message
}

func NewAddDataMenu(sub *subscribtion.Subscription, client *grpc_client.Client, pages *tview.Pages, message *Message) *AddDataMenu {
	addDataMenu := &AddDataMenu{
		sub:     sub,
		client:  client,
		pages:   pages,
		message: message,
	}
	sub.SubscribeEvent(events.ShowAddDataMenu, addDataMenu.GetNotifier())
	return addDataMenu
}

// Show Меню для выбора типа данных при добавлении
func (adm *AddDataMenu) Show() {
	menu := tview.NewList().
		AddItem("Пара логин/пароль", "", '1', func() {
			adm.sub.NotifyEvent(events.ShowAddDataForm, events.AddDataFormEvent{Type: "login_password"})
		}).
		AddItem("Текстовые данные", "", '2', func() {
			adm.sub.NotifyEvent(events.ShowAddDataForm, events.AddDataFormEvent{Type: "text"})
		}).
		AddItem("Бинарные данные", "", '3', func() {
			adm.sub.NotifyEvent(events.ShowAddDataForm, events.AddDataFormEvent{Type: "binary"})
		}).
		AddItem("Данные банковской карты", "", '4', func() {
			adm.sub.NotifyEvent(events.ShowAddDataForm, events.AddDataFormEvent{Type: "bank_card"})
		}).
		AddItem("Назад", "", 'b', func() {
			adm.sub.NotifyEvent(events.ShowMainMenu, nil)
		})
	menu.SetTitle("Выберите тип данных для добавления").SetBorder(true)

	adm.pages.AddAndSwitchToPage("addDataMenu", menu, true)
}

func (adm *AddDataMenu) GetNotifier() subscribtion.Notifier {
	return func(_ any) {
		adm.Show()
	}
}
