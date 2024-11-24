package tui

import (
	"github.com/moonicy/goph-keeper-yandex/internal/events"
	"github.com/moonicy/goph-keeper-yandex/internal/grpc_client"
	"github.com/moonicy/goph-keeper-yandex/internal/subscribtion"
	"github.com/rivo/tview"
)

type MainMenu struct {
	sub    *subscribtion.Subscription
	client *grpc_client.Client
	pages  *tview.Pages
}

func NewMainMenu(sub *subscribtion.Subscription, client *grpc_client.Client, pages *tview.Pages) *MainMenu {
	mainMenu := &MainMenu{sub: sub, client: client, pages: pages}
	sub.SubscribeEvent(events.ShowMainMenu, mainMenu.GetNotifier())
	return mainMenu
}

// Show Функция для отображения главного меню
func (mm *MainMenu) Show() {
	menu := tview.NewList().
		AddItem("Добавить данные", "", 'a', func() {
			mm.sub.NotifyEvent(events.ShowAddDataMenu, nil)
		}).
		AddItem("Обновить данные", "", 'u', func() {
			mm.sub.NotifyEvent(events.ShowUpdateDataForm, nil)
		}).
		AddItem("Получить данные", "", 'g', func() {
			mm.sub.NotifyEvent(events.ShowGetData, nil)
		}).
		AddItem("Удалить данные", "", 'r', func() {
			mm.sub.NotifyEvent(events.ShowRemoveDataForm, nil)
		}).
		AddItem("Выйти", "", 'l', func() {
			mm.client.Logout()
			mm.sub.NotifyEvent(events.ShowLoginForm, nil)
		})
	menu.SetTitle("Главное меню").SetBorder(true)

	mm.pages.AddAndSwitchToPage("mainmenu", menu, true)
}

func (mm *MainMenu) GetNotifier() subscribtion.Notifier {
	return func(_ any) {
		mm.Show()
	}
}
