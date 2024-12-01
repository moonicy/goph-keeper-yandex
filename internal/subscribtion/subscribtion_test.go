package subscribtion

import (
	"testing"

	evt "github.com/moonicy/goph-keeper-yandex/internal/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockNotifier struct {
	mock.Mock
}

func (m *mockNotifier) Notify(eventData any) {
	m.Called(eventData)
}

func TestNewSubscription(t *testing.T) {
	sub := NewSubscription()
	assert.NotNil(t, sub)
	assert.Empty(t, sub.events)
}

func TestSubscription_SubscribeEvent(t *testing.T) {
	sub := NewSubscription()

	event := evt.Event("TestEvent")
	notifier := func(data any) {}

	sub.SubscribeEvent(event, notifier)

	assert.Contains(t, sub.events, event)
	assert.NotNil(t, sub.events[event])
}

func TestSubscription_NotifyEvent_Success(t *testing.T) {
	sub := NewSubscription()

	event := evt.Event("TestEvent")
	mockNotifier := &mockNotifier{}
	mockNotifier.On("Notify", "TestData").Return()

	sub.SubscribeEvent(event, mockNotifier.Notify)

	sub.NotifyEvent(event, "TestData")

	mockNotifier.AssertCalled(t, "Notify", "TestData")
}
