package tui

import (
	"encoding/json"
	"fmt"
	"github.com/moonicy/goph-keeper-yandex/internal/events"
	"github.com/moonicy/goph-keeper-yandex/internal/grpc_client"
	"github.com/moonicy/goph-keeper-yandex/internal/subscribtion"
	"github.com/rivo/tview"
)

type RemoveDataForm struct {
	sub     *subscribtion.Subscription
	client  *grpc_client.Client
	pages   *tview.Pages
	message *Message
}

func NewRemoveDataForm(sub *subscribtion.Subscription, client *grpc_client.Client, pages *tview.Pages, message *Message) *RemoveDataForm {
	removeDataForm := &RemoveDataForm{
		sub:     sub,
		client:  client,
		pages:   pages,
		message: message,
	}
	sub.SubscribeEvent(events.ShowRemoveDataForm, removeDataForm.GetNotifier())
	return removeDataForm
}

// Show Функция для отображения списка данных для удаления
func (rdf *RemoveDataForm) Show() {
	data, err := rdf.client.GetData()
	if err != nil {
		rdf.message.ShowError(err)
		return
	}

	list := tview.NewList()
	for _, d := range data {
		var dataStruct map[string]interface{}
		err = json.Unmarshal(d.Data, &dataStruct)
		if err != nil {
			rdf.message.ShowError(fmt.Errorf("ошибка десериализации данных: %v", err))
			return
		}

		dataType, _ := dataStruct["type"].(string)
		metadata, _ := dataStruct["metadata"].(string)

		itemText := fmt.Sprintf("ID: %d | Тип: %s | Мета: %s", d.ID, dataType, metadata)
		idCopy := d.ID
		itemTextCopy := itemText
		list.AddItem(itemTextCopy, "", 0, func() {
			// Подтверждение удаления
			confirmModal := tview.NewModal().
				SetText(fmt.Sprintf("Вы уверены, что хотите удалить запись?\n\n%s", itemTextCopy)).
				AddButtons([]string{"Да", "Нет"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					rdf.pages.RemovePage("confirmModal")
					if buttonLabel == "Да" {
						err = rdf.client.RemoveData(idCopy)
						if err != nil {
							rdf.message.ShowError(err)
							return
						}
						rdf.message.ShowMessage("Данные успешно удалены")
						rdf.sub.NotifyEvent(events.ShowMainMenu, nil)
					} else {
						// Возвращаемся к списку данных для удаления
						rdf.pages.SwitchToPage("removeData")
					}
				})
			confirmModal.SetTitle("Подтверждение удаления").SetBorder(true)
			rdf.pages.AddAndSwitchToPage("confirmModal", confirmModal, true)
		})
	}
	list.AddItem("Назад", "", 'b', func() {
		rdf.sub.NotifyEvent(events.ShowMainMenu, nil)
	})
	list.SetTitle("Выберите данные для удаления").SetBorder(true)

	rdf.pages.AddAndSwitchToPage("removeData", list, true)
}

func (rdf *RemoveDataForm) GetNotifier() subscribtion.Notifier {
	return func(_ any) {
		rdf.Show()
	}
}
