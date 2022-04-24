package eventutils

import (
	"context"

	"github.com/EventStore/EventStore-Client-Go/esdb"
)

type EventHandler struct {
	subscription *esdb.PersistentSubscription
	handlers     map[string]func(rawEvent *esdb.SubscriptionEvent)
}

func NewEventHandler(client *esdb.Client, groupName string) EventHandler {
	sub, err := client.ConnectToPersistentSubscriptionToAll(
		context.Background(),
		groupName,
		esdb.ConnectToPersistentSubscriptionOptions{},
	)
	if err != nil {
		panic(err)
	}
	return EventHandler{
		subscription: sub,
		handlers:     make(map[string]func(rawEvent *esdb.SubscriptionEvent)),
	}
}

func (handler EventHandler) HandleEvent(eventName string, fn func(rawEvent *esdb.SubscriptionEvent)) EventHandler {
	handler.handlers[eventName] = fn
	return handler
}

func (handler EventHandler) Start() {
	go func() {
		for {
			event := handler.subscription.Recv()
			if event.EventAppeared != nil {
				handler.handlers[event.EventAppeared.Event.EventType](event)
				handler.subscription.Ack(event.EventAppeared)
			}
			if event.SubscriptionDropped != nil {
				break
			}
		}
	}()
}
