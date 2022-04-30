package internal

import (
	"encoding/json"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/Resta-Inc/resta/menu/commands/internal/entities"
	"github.com/Resta-Inc/resta/pkg/events"
	"github.com/Resta-Inc/resta/pkg/eventutils"
)

type EventHandler struct {
	entityRepository eventutils.IEntityRepository
}

func NewEventHandler(repo eventutils.IEntityRepository) EventHandler {
	return EventHandler{
		entityRepository: repo,
	}
}

func (eventHandler EventHandler) HandleCategoryCreated(rawEvent *esdb.SubscriptionEvent) error {
	_, rawData := eventutils.GetRawDataFromSerializedEvent(rawEvent.EventAppeared.Event.Data)
	var event events.CategoryCreated
	err := json.Unmarshal(rawData, &event)
	if err != nil {
		return err
	}
	menu, err := eventHandler.entityRepository.GetEntity(entities.EmptyMenu(), event.ParentMenuID)
	if err != nil {
		return err
	}
	menu = menu.(entities.Menu).AddCategory(event.GetEntityID())
	err = eventHandler.entityRepository.SaveEntity(menu)
	return err
}
