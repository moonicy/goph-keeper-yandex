package main

import (
	_ "embed"
	"github.com/moonicy/goph-keeper-yandex/internal/config"
	"github.com/moonicy/goph-keeper-yandex/internal/grpc_client"
	"github.com/moonicy/goph-keeper-yandex/internal/subscribtion"
	"github.com/moonicy/goph-keeper-yandex/internal/tui"
	"log"

	"github.com/rivo/tview"
)

func main() {
	cfg := config.NewClientConfig()

	// Устанавливаем соединение с gRPC сервером
	client, err := grpc_client.NewClient(cfg.Host)
	if err != nil {
		log.Fatalf("Не удалось подключиться: %v", err)
	}
	defer client.Close()

	app := tview.NewApplication()
	pages := tview.NewPages()

	app.SetRoot(pages, true)

	sub := subscribtion.NewSubscription()

	message := tui.NewMessage(pages)
	_ = tui.NewMainMenu(sub, client, pages)
	loginForm := tui.NewLoginForm(sub, client, pages, message)
	_ = tui.NewRegisterForm(sub, client, pages, message)
	_ = tui.NewAddDataMenu(sub, client, pages, message)
	_ = tui.NewAddDataForm(sub, client, pages, message)
	_ = tui.NewUpdateDataForm(sub, client, pages, message)
	_ = tui.NewEditDataForm(sub, client, pages, message)
	_ = tui.NewGetData(sub, client, pages, message)
	_ = tui.NewDataDetails(sub, client, pages, message)
	_ = tui.NewRemoveDataForm(sub, client, pages, message)

	loginForm.Show()

	if err = app.Run(); err != nil {
		log.Fatalf("Ошибка запуска приложения: %v", err)
	}
}
