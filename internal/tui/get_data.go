package tui

import (
	"encoding/json"
	"fmt"
	"github.com/moonicy/goph-keeper-yandex/internal/events"
	"github.com/moonicy/goph-keeper-yandex/internal/grpc_client"
	"github.com/moonicy/goph-keeper-yandex/internal/subscribtion"
	"github.com/rivo/tview"
)

type GetData struct {
	sub     *subscribtion.Subscription
	client  *grpc_client.Client
	pages   *tview.Pages
	message *Message
}

func NewGetData(sub *subscribtion.Subscription, client *grpc_client.Client, pages *tview.Pages, message *Message) *GetData {
	getData := &GetData{
		sub:     sub,
		client:  client,
		pages:   pages,
		message: message,
	}
	sub.SubscribeEvent(events.ShowGetData, getData.GetNotifier())
	return getData
}

// Show Функция для отображения списка данных
func (gd *GetData) Show() {
	data, err := gd.client.GetData()
	if err != nil {
		gd.message.ShowError(err)
		return
	}

	list := tview.NewList()
	for _, d := range data {
		var dataStruct map[string]interface{}
		err = json.Unmarshal(d.Data, &dataStruct)
		if err != nil {
			gd.message.ShowError(fmt.Errorf("ошибка десериализации данных: %v", err))
			return
		}

		dataType, _ := dataStruct["type"].(string)
		metadata, _ := dataStruct["metadata"].(string)

		var displayData string
		switch dataType {
		case "login_password":
			dataMap, _ := dataStruct["data"].(map[string]interface{})
			login, _ := dataMap["login"].(string)
			password, _ := dataMap["password"].(string)
			displayData = fmt.Sprintf("Логин: %s\nПароль: %s", login, password)
		case "text":
			text, _ := dataStruct["data"].(string)
			displayData = fmt.Sprintf("Текст:\n%s", text)
		case "binary":
			displayData = "Бинарные данные (нельзя отобразить)"
		case "bank_card":
			dataMap, _ := dataStruct["data"].(map[string]interface{})
			number, _ := dataMap["number"].(string)
			expiry, _ := dataMap["expiry"].(string)
			cvv, _ := dataMap["cvv"].(string)
			holder, _ := dataMap["holder"].(string)
			displayData = fmt.Sprintf("Номер карты: %s\nСрок действия: %s\nCVV: %s\nИмя держателя: %s",
				number, expiry, cvv, holder)
		default:
			displayData = "Неизвестный тип данных"
		}

		itemText := fmt.Sprintf("ID: %d | Тип: %s | Мета: %s", d.ID, dataType, metadata)
		idCopy := d.ID
		dataTypeCopy := dataType
		metadataCopy := metadata
		displayDataCopy := displayData
		list.AddItem(itemText, "", 0, func() {
			gd.sub.NotifyEvent(events.ShowDataDetails, events.DataDetailsFormEvent{
				ID:          idCopy,
				DataType:    dataTypeCopy,
				Metadata:    metadataCopy,
				DisplayData: displayDataCopy,
			})
		})
	}
	list.AddItem("Назад", "", 'b', func() {
		gd.sub.NotifyEvent(events.ShowMainMenu, nil)
	})
	list.SetTitle("Список данных").SetBorder(true)

	gd.pages.AddAndSwitchToPage("getData", list, true)
}

func (gd *GetData) GetNotifier() subscribtion.Notifier {
	return func(_ any) {
		gd.Show()
	}
}
