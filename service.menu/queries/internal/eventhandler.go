package internal

import (
	"encoding/json"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/Resta-Inc/resta/pkg/events"
	"github.com/Resta-Inc/resta/pkg/eventutils"
)

type MenuEventHandler struct {
	menuRepository IMenuRepository
}

func NewMenuEventHandler(repo IMenuRepository) MenuEventHandler {
	return MenuEventHandler{
		menuRepository: repo,
	}
}

func (menuEventHandler MenuEventHandler) HandleMenuCreated(rawEvent *esdb.SubscriptionEvent) error {
	_, rawData := eventutils.GetRawDataFromSerializedEvent(rawEvent.EventAppeared.Event.Data)
	var event events.MenuCreated
	err := json.Unmarshal(rawData, &event)
	if err != nil {
		panic(err)
	}
	err = menuEventHandler.menuRepository.CreateMenu(event.GetEntityID(), event.Name)
	return err
}
