package tui

import (
	"encoding/json"
	"fmt"
	"github.com/moonicy/goph-keeper-yandex/internal/events"
	"github.com/moonicy/goph-keeper-yandex/internal/grpc_client"
	"github.com/moonicy/goph-keeper-yandex/internal/subscribtion"
	"github.com/rivo/tview"
	"io"
	"log"
	"strings"
)

type AddDataForm struct {
	sub     *subscribtion.Subscription
	client  *grpc_client.Client
	pages   *tview.Pages
	message *Message
}

func NewAddDataForm(sub *subscribtion.Subscription, client *grpc_client.Client, pages *tview.Pages, message *Message) *AddDataForm {
	addDataForm := &AddDataForm{
		sub:     sub,
		client:  client,
		pages:   pages,
		message: message,
	}
	sub.SubscribeEvent(events.ShowAddDataForm, addDataForm.GetNotifier())
	return addDataForm
}

// Show Функция для отображения формы добавления данных
func (adf *AddDataForm) Show(dataType string) {
	var addDataForm *tview.Form

	addDataForm = tview.NewForm()
	// В зависимости от типа данных добавляем соответствующие поля
	switch dataType {
	case loginPasswordFormType:
		addDataForm.
			AddInputField("Логин", "", 50, nil, nil).
			AddPasswordField("Пароль", "", 50, '*', nil)
	case textFormType:
		addDataForm.
			AddInputField("Текст", "", 50, nil, nil)
	case binaryFormType:
		addDataForm.
			AddInputField("Путь к файлу", "", 50, nil, nil)
	case bankCardFormType:
		addDataForm.
			AddInputField("Номер карты", "", 20, nil, nil).
			AddInputField("Срок действия (MM/YY)", "", 5, nil, nil).
			AddInputField("CVV", "", 3, nil, nil).
			AddInputField("Имя держателя", "", 50, nil, nil)
	}

	addDataForm.
		AddInputField("Метаинформация", "", 100, nil, nil).
		AddButton("Добавить", func() {
			var dataBytes []byte
			var err error

			// Получаем метаинформацию
			metadata := addDataForm.GetFormItemByLabel("Метаинформация").(*tview.InputField).GetText()

			dataStruct := map[string]interface{}{
				"type":     dataType,
				"metadata": metadata,
				"data":     nil,
			}

			// Формируем данные в зависимости от типа
			switch dataType {
			case loginPasswordFormType:
				login := addDataForm.GetFormItemByLabel("Логин").(*tview.InputField).GetText()
				password := addDataForm.GetFormItemByLabel("Пароль").(*tview.InputField).GetText()
				dataStruct["data"] = map[string]string{
					"login":    login,
					"password": password,
				}
			case textFormType:
				text := addDataForm.GetFormItemByLabel("Текст").(*tview.InputField).GetText()
				dataStruct["data"] = text
			case binaryFormType:
				filePath := addDataForm.GetFormItemByLabel("Путь к файлу").(*tview.InputField).GetText()
				dataBytes, err = readFileAsBytes(filePath)
				if err != nil {
					adf.message.ShowError(fmt.Errorf("ошибка чтения файла: %v", err))
					return
				}
				// Преобразуем бинарные данные в строку Base64 для передачи в JSON
				dataStruct["data"] = dataBytes
			case bankCardFormType:
				cardNumber := addDataForm.GetFormItemByLabel("Номер карты").(*tview.InputField).GetText()
				expiryDate := addDataForm.GetFormItemByLabel("Срок действия (MM/YY)").(*tview.InputField).GetText()
				cvv := addDataForm.GetFormItemByLabel("CVV").(*tview.InputField).GetText()
				cardHolder := addDataForm.GetFormItemByLabel("Имя держателя").(*tview.InputField).GetText()
				dataStruct["data"] = map[string]string{
					"number": cardNumber,
					"expiry": expiryDate,
					"cvv":    cvv,
					"holder": cardHolder,
				}
			}

			dataBytes, err = json.Marshal(dataStruct)
			if err != nil {
				adf.message.ShowError(fmt.Errorf("ошибка сериализации данных: %v", err))
				return
			}

			err = adf.client.AddData(dataBytes)
			if err != nil {
				adf.message.ShowError(err)
				return
			}

			adf.message.ShowMessage("Данные успешно добавлены")
			adf.sub.NotifyEvent(events.ShowMainMenu, nil)
		}).
		AddButton("Назад", func() {
			adf.sub.NotifyEvent(events.ShowMainMenu, nil)
		})
	addDataForm.SetTitle("Добавить данные").SetBorder(true)

	adf.pages.AddAndSwitchToPage("addData", addDataForm, true)
}

func (adf *AddDataForm) GetNotifier() subscribtion.Notifier {
	return func(eventData any) {
		data, ok := eventData.(events.AddDataFormEvent)
		if !ok {
			log.Fatal("Wrong event data type")
		}
		adf.Show(data.Type)
	}
}

// Функция для чтения файла в виде байтов
func readFileAsBytes(filePath string) ([]byte, error) {
	data, err := io.ReadAll(strings.NewReader(filePath))
	if err != nil {
		return nil, err
	}
	return data, nil
}
