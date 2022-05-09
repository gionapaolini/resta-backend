package internal

import (
	"encoding/json"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/Resta-Inc/resta/menu/commands/internal/entities"
	"github.com/Resta-Inc/resta/pkg/events"
	"github.com/Resta-Inc/resta/pkg/eventutils"
)

type MenuEventHandler struct {
	entityRepository eventutils.IEntityRepository
}

func NewMenuEventHandler(repo eventutils.IEntityRepository) MenuEventHandler {
	return MenuEventHandler{
		entityRepository: repo,
	}
}

func (eventHandler MenuEventHandler) HandleCategoryCreated(rawEvent *esdb.SubscriptionEvent) error {
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

func (eventHandler MenuEventHandler) HandleSubCategoryCreated(rawEvent *esdb.SubscriptionEvent) error {
	_, rawData := eventutils.GetRawDataFromSerializedEvent(rawEvent.EventAppeared.Event.Data)
	var event events.SubCategoryCreated
	err := json.Unmarshal(rawData, &event)
	if err != nil {
		return err
	}
	category, err := eventHandler.entityRepository.GetEntity(entities.EmptyCategory(), event.ParentCategoryID)
	if err != nil {
		return err
	}
	category = category.(entities.Category).AddSubCategory(event.GetEntityID())
	err = eventHandler.entityRepository.SaveEntity(category)
	return err
}
