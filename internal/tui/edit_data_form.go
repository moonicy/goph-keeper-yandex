package tui

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/moonicy/goph-keeper-yandex/internal/events"
	"github.com/moonicy/goph-keeper-yandex/internal/grpc_client"
	"github.com/moonicy/goph-keeper-yandex/internal/subscribtion"
	"github.com/rivo/tview"
	"log"
)

type EditDataForm struct {
	sub     *subscribtion.Subscription
	client  *grpc_client.Client
	pages   *tview.Pages
	message *Message
}

func NewEditDataForm(sub *subscribtion.Subscription, client *grpc_client.Client, pages *tview.Pages, message *Message) *EditDataForm {
	editDataForm := &EditDataForm{
		sub:     sub,
		client:  client,
		pages:   pages,
		message: message,
	}
	sub.SubscribeEvent(events.ShowEditDataForm, editDataForm.GetNotifier())
	return editDataForm
}

// Show Функция для отображения формы редактирования данных
func (edf *EditDataForm) Show(id uint64, dataStruct map[string]interface{}) {
	var editDataForm *tview.Form

	dataType, _ := dataStruct["type"].(string)
	metadata, _ := dataStruct["metadata"].(string)
	editDataForm = tview.NewForm()

	// В зависимости от типа данных добавляем соответствующие поля и предзаполняем их текущими значениями
	switch dataType {
	case "login_password":
		dataMap, _ := dataStruct["data"].(map[string]interface{})
		login, _ := dataMap["login"].(string)
		password, _ := dataMap["password"].(string)
		editDataForm.
			AddInputField("Логин", login, 50, nil, nil).
			AddPasswordField("Пароль", password, 50, '*', nil)
	case "text":
		text, _ := dataStruct["data"].(string)
		editDataForm.
			AddInputField("Текст", text, 50, nil, nil)
	case "binary":
		// Бинарные данные нельзя отобразить; предлагаем выбрать новый файл
		editDataForm.
			AddInputField("Путь к новому файлу (оставьте пустым, чтобы оставить прежние данные)", "", 50, nil, nil)
	case "bank_card":
		dataMap, _ := dataStruct["data"].(map[string]interface{})
		number, _ := dataMap["number"].(string)
		expiry, _ := dataMap["expiry"].(string)
		cvv, _ := dataMap["cvv"].(string)
		holder, _ := dataMap["holder"].(string)
		editDataForm.
			AddInputField("Номер карты", number, 20, nil, nil).
			AddInputField("Срок действия (MM/YY)", expiry, 5, nil, nil).
			AddInputField("CVV", cvv, 3, nil, nil).
			AddInputField("Имя держателя", holder, 50, nil, nil)
	default:
		edf.message.ShowError(errors.New("неизвестный тип данных"))
		return
	}

	editDataForm.
		AddInputField("Метаинформация", metadata, 100, nil, nil).
		AddButton("Сохранить", func() {
			var dataBytes []byte
			var err error

			newMetadata := editDataForm.GetFormItemByLabel("Метаинформация").(*tview.InputField).GetText()

			newDataStruct := map[string]interface{}{
				"type":     dataType,
				"metadata": newMetadata,
				"data":     nil,
			}

			// Формируем данные в зависимости от типа
			switch dataType {
			case "login_password":
				newLogin := editDataForm.GetFormItemByLabel("Логин").(*tview.InputField).GetText()
				newPassword := editDataForm.GetFormItemByLabel("Пароль").(*tview.InputField).GetText()
				newDataStruct["data"] = map[string]string{
					"login":    newLogin,
					"password": newPassword,
				}
			case "text":
				newText := editDataForm.GetFormItemByLabel("Текст").(*tview.InputField).GetText()
				newDataStruct["data"] = newText
			case "binary":
				filePath := editDataForm.GetFormItemByLabel("Путь к новому файлу (оставьте пустым, чтобы оставить прежние данные)").(*tview.InputField).GetText()
				if filePath != "" {
					dataBytes, err = readFileAsBytes(filePath)
					if err != nil {
						edf.message.ShowError(fmt.Errorf("ошибка чтения файла: %v", err))
						return
					}
					newDataStruct["data"] = dataBytes
				} else {
					// Оставляем прежние данные
					newDataStruct["data"] = dataStruct["data"]
				}
			case "bank_card":
				newNumber := editDataForm.GetFormItemByLabel("Номер карты").(*tview.InputField).GetText()
				newExpiry := editDataForm.GetFormItemByLabel("Срок действия (MM/YY)").(*tview.InputField).GetText()
				newCVV := editDataForm.GetFormItemByLabel("CVV").(*tview.InputField).GetText()
				newHolder := editDataForm.GetFormItemByLabel("Имя держателя").(*tview.InputField).GetText()
				newDataStruct["data"] = map[string]string{
					"number": newNumber,
					"expiry": newExpiry,
					"cvv":    newCVV,
					"holder": newHolder,
				}
			default:
				edf.message.ShowError(errors.New("неизвестный тип данных"))
				return
			}

			dataBytes, err = json.Marshal(newDataStruct)
			if err != nil {
				edf.message.ShowError(fmt.Errorf("ошибка сериализации данных: %v", err))
				return
			}

			err = edf.client.UpdateData(id, dataBytes)
			if err != nil {
				edf.message.ShowError(err)
				return
			}

			edf.message.ShowMessage("Данные успешно обновлены")
			edf.sub.NotifyEvent(events.ShowMainMenu, nil)
		}).
		AddButton("Отмена", func() {
			edf.sub.NotifyEvent(events.ShowMainMenu, nil)
		})
	editDataForm.SetTitle("Редактировать данные").SetBorder(true)

	edf.pages.AddAndSwitchToPage("editData", editDataForm, true)
}

func (edf *EditDataForm) GetNotifier() subscribtion.Notifier {
	return func(eventData any) {
		data, ok := eventData.(events.EditDataFormEvent)
		if !ok {
			log.Fatal("Wrong event data type")
		}
		edf.Show(data.ID, data.DataStruct)
	}
}
