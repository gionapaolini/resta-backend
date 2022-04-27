package eventutils

import (
	"context"

	"github.com/EventStore/EventStore-Client-Go/esdb"
)

type EventHandler struct {
	subscription *esdb.PersistentSubscription
	handlers     map[string]func(rawEvent *esdb.SubscriptionEvent) error
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
		handlers:     make(map[string]func(rawEvent *esdb.SubscriptionEvent) error),
	}
}

func (handler EventHandler) HandleEvent(eventName string, fn func(rawEvent *esdb.SubscriptionEvent) error) EventHandler {
	handler.handlers[eventName] = fn
	return handler
}

func (handler EventHandler) Start() {
	go func() {
		for {
			event := handler.subscription.Recv()
			if event.EventAppeared != nil {
				err := handler.handlers[event.EventAppeared.Event.EventType](event)
				if err != nil {
					handler.subscription.Nack(err.Error(), esdb.Nack_Retry)
				} else {
					handler.subscription.Ack(event.EventAppeared)
				}
			}
			if event.SubscriptionDropped != nil {
				break
			}
		}
	}()
}
