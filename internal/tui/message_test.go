package tui

import (
	"errors"
	"testing"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestNewMessage(t *testing.T) {
	pages := tview.NewPages()
	msg := NewMessage(pages)

	assert.NotNil(t, msg)
	assert.Equal(t, pages, msg.pages)
}

func TestMessage_ShowError(t *testing.T) {
	pages := tview.NewPages()
	msg := NewMessage(pages)

	err := errors.New("test error message")
	msg.ShowError(err)

	hasPage := pages.HasPage("modal")
	assert.True(t, hasPage)

	msg.pages.RemovePage("modal")

	hasPage = pages.HasPage("modal")
	assert.False(t, hasPage)
}

func TestMessage_ShowMessage(t *testing.T) {
	pages := tview.NewPages()
	msg := NewMessage(pages)

	message := "Test message"
	msg.ShowMessage(message)

	hasPage := pages.HasPage("modal")
	assert.True(t, hasPage)

	msg.pages.RemovePage("modal")

	hasPage = pages.HasPage("modal")
	assert.False(t, hasPage)
}
