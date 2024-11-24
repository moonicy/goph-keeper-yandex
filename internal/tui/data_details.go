package tui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/moonicy/goph-keeper-yandex/internal/events"
	"github.com/moonicy/goph-keeper-yandex/internal/grpc_client"
	"github.com/moonicy/goph-keeper-yandex/internal/subscribtion"
	"github.com/rivo/tview"
	"log"
)

type DataDetails struct {
	sub     *subscribtion.Subscription
	client  *grpc_client.Client
	pages   *tview.Pages
	message *Message
}

func NewDataDetails(sub *subscribtion.Subscription, client *grpc_client.Client, pages *tview.Pages, message *Message) *DataDetails {
	dataDetails := &DataDetails{
		sub:     sub,
		client:  client,
		pages:   pages,
		message: message,
	}
	sub.SubscribeEvent(events.ShowDataDetails, dataDetails.GetNotifier())
	return dataDetails
}

// Show Функция для отображения деталей данных
func (dd *DataDetails) Show(id uint64, dataType, metadata, content string) {
	textView := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true).
		SetWrap(true).
		SetScrollable(true).
		SetText(fmt.Sprintf("[yellow]ID: %d\nТип: %s\nМета: %s\n\n[white]%s", id, dataType, metadata, content))
	textView.SetBorder(true).SetTitle("Детали данных")

	// Устанавливаем обработчик ввода для перехвата нажатий клавиш
	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc, tcell.KeyEnter, tcell.KeyBackspace, tcell.KeyBackspace2:
			// Переключаемся на страницу со списком данных
			dd.pages.SwitchToPage("getData")
			dd.pages.RemovePage("dataDetails")
			return nil
		case tcell.KeyRune:
			if event.Rune() == ' ' {
				dd.pages.SwitchToPage("getData")
				dd.pages.RemovePage("dataDetails")
				return nil
			}
		}
		return event
	})

	// Устанавливаем функцию завершения при нажатии Esc
	textView.SetDoneFunc(func(key tcell.Key) {
		// Переключаемся на страницу со списком данных
		dd.pages.SwitchToPage("getData")
		dd.pages.RemovePage("dataDetails")
	})

	// Добавляем TextView как новую страницу в менеджер страниц
	dd.pages.AddAndSwitchToPage("dataDetails", textView, true)
}

func (dd *DataDetails) GetNotifier() subscribtion.Notifier {
	return func(eventData any) {
		data, ok := eventData.(events.DataDetailsFormEvent)
		if !ok {
			log.Fatal("Wrong event data type")
		}
		dd.Show(data.ID, data.DataType, data.Metadata, data.DisplayData)
	}
}
