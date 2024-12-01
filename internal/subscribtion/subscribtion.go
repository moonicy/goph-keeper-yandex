package subscribtion

import (
	evt "github.com/moonicy/goph-keeper-yandex/internal/events"
	"log"
)

type Notifier func(eventData any)

type Subscription struct {
	events map[evt.Event]Notifier
}

func NewSubscription() *Subscription {
	return &Subscription{
		events: make(map[evt.Event]Notifier),
	}
}

func (sub *Subscription) SubscribeEvent(event evt.Event, action Notifier) {
	sub.events[event] = action
}

func (sub *Subscription) NotifyEvent(event evt.Event, eventData any) {
	f, ok := sub.events[event]
	if !ok {
		log.Fatalf("key %s not found", event)
	}
	f(eventData)
}
