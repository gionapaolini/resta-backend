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
		return err
	}
	err = menuEventHandler.menuRepository.CreateMenu(event.GetEntityID(), event.Name)
	return err
}

func (menuEventHandler MenuEventHandler) HandleMenuEnabled(rawEvent *esdb.SubscriptionEvent) error {
	_, rawData := eventutils.GetRawDataFromSerializedEvent(rawEvent.EventAppeared.Event.Data)
	var event events.MenuEnabled
	err := json.Unmarshal(rawData, &event)
	if err != nil {
		return err
	}
	err = menuEventHandler.menuRepository.EnableMenu(event.GetEntityID())
	return err
}

func (menuEventHandler MenuEventHandler) HandleMenuDisabled(rawEvent *esdb.SubscriptionEvent) error {
	_, rawData := eventutils.GetRawDataFromSerializedEvent(rawEvent.EventAppeared.Event.Data)
	var event events.MenuDisabled
	err := json.Unmarshal(rawData, &event)
	if err != nil {
		return err
	}
	err = menuEventHandler.menuRepository.DisableMenu(event.GetEntityID())
	return err
}

func (menuEventHandler MenuEventHandler) HandleMenuNameChanged(rawEvent *esdb.SubscriptionEvent) error {
	_, rawData := eventutils.GetRawDataFromSerializedEvent(rawEvent.EventAppeared.Event.Data)
	var event events.MenuNameChanged
	err := json.Unmarshal(rawData, &event)
	if err != nil {
		return err
	}
	err = menuEventHandler.menuRepository.ChangeMenuName(event.GetEntityID(), event.NewName)
	return err
}
