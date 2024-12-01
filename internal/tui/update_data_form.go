package tui

import (
	"encoding/json"
	"fmt"
	"github.com/moonicy/goph-keeper-yandex/internal/events"
	"github.com/moonicy/goph-keeper-yandex/internal/grpc_client"
	"github.com/moonicy/goph-keeper-yandex/internal/subscribtion"
	"github.com/rivo/tview"
)

type UpdateDataForm struct {
	sub     *subscribtion.Subscription
	client  *grpc_client.Client
	pages   *tview.Pages
	message *Message
}

func NewUpdateDataForm(sub *subscribtion.Subscription, client *grpc_client.Client, pages *tview.Pages, message *Message) *UpdateDataForm {
	updateDataForm := &UpdateDataForm{
		sub:     sub,
		client:  client,
		pages:   pages,
		message: message,
	}
	sub.SubscribeEvent(events.ShowUpdateDataForm, updateDataForm.GetNotifier())
	return updateDataForm
}

// Show Функция для отображения списка данных для редактирования
func (udf *UpdateDataForm) Show() {
	data, err := udf.client.GetData()
	if err != nil {
		udf.message.ShowError(err)
		return
	}

	list := tview.NewList()
	for _, d := range data {
		var dataStruct map[string]interface{}
		err = json.Unmarshal(d.Data, &dataStruct)
		if err != nil {
			udf.message.ShowError(fmt.Errorf("ошибка десериализации данных: %v", err))
			return
		}

		dataType, _ := dataStruct["type"].(string)
		metadata, _ := dataStruct["metadata"].(string)

		itemText := fmt.Sprintf("ID: %d | Тип: %s | Мета: %s", d.ID, dataType, metadata)
		dataCopy := d
		dataStructCopy := dataStruct
		list.AddItem(itemText, "", 0, func() {
			udf.sub.NotifyEvent(events.ShowEditDataForm, events.EditDataFormEvent{
				ID:         dataCopy.ID,
				DataStruct: dataStructCopy,
			})
		})
	}
	list.AddItem("Назад", "", 'b', func() {
		udf.sub.NotifyEvent(events.ShowMainMenu, nil)
	})
	list.SetTitle("Выберите данные для редактирования").SetBorder(true)

	udf.pages.AddAndSwitchToPage("updateData", list, true)
}

func (udf *UpdateDataForm) GetNotifier() subscribtion.Notifier {
	return func(_ any) {
		udf.Show()
	}
}
