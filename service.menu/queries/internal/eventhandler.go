package internal

import (
	"encoding/json"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/Resta-Inc/resta/pkg/events"
	"github.com/Resta-Inc/resta/pkg/eventutils2"
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
	recordedEvent := eventutils2.DeserializeRecordedEvent(rawEvent.EventAppeared.Event)
	var event events.MenuCreated
	err := json.Unmarshal(recordedEvent.Data, &event)
	if err != nil {
		return err
	}
	err = menuEventHandler.menuRepository.CreateMenu(event.GetEntityID(), event.Name)
	return err
}

func (menuEventHandler MenuEventHandler) HandleMenuEnabled(rawEvent *esdb.SubscriptionEvent) error {
	recordedEvent := eventutils2.DeserializeRecordedEvent(rawEvent.EventAppeared.Event)
	var event events.MenuEnabled
	err := json.Unmarshal(recordedEvent.Data, &event)
	if err != nil {
		return err
	}
	err = menuEventHandler.menuRepository.EnableMenu(event.GetEntityID())
	return err
}

func (menuEventHandler MenuEventHandler) HandleMenuDisabled(rawEvent *esdb.SubscriptionEvent) error {
	recordedEvent := eventutils2.DeserializeRecordedEvent(rawEvent.EventAppeared.Event)
	var event events.MenuDisabled
	err := json.Unmarshal(recordedEvent.Data, &event)
	if err != nil {
		return err
	}
	err = menuEventHandler.menuRepository.DisableMenu(event.GetEntityID())
	return err
}

func (menuEventHandler MenuEventHandler) HandleMenuNameChanged(rawEvent *esdb.SubscriptionEvent) error {
	recordedEvent := eventutils2.DeserializeRecordedEvent(rawEvent.EventAppeared.Event)
	var event events.MenuNameChanged
	err := json.Unmarshal(recordedEvent.Data, &event)
	if err != nil {
		return err
	}
	err = menuEventHandler.menuRepository.ChangeMenuName(event.GetEntityID(), event.NewName)
	return err
}

func (menuEventHandler MenuEventHandler) HandleCategoryCreated(rawEvent *esdb.SubscriptionEvent) error {
	recordedEvent := eventutils2.DeserializeRecordedEvent(rawEvent.EventAppeared.Event)
	var event events.CategoryCreated
	err := json.Unmarshal(recordedEvent.Data, &event)
	if err != nil {
		return err
	}
	err = menuEventHandler.menuRepository.CreateCategory(event.GetEntityID(), event.Name)
	return err
}

func (menuEventHandler MenuEventHandler) HandleCategoryAddedToMenu(rawEvent *esdb.SubscriptionEvent) error {
	recordedEvent := eventutils2.DeserializeRecordedEvent(rawEvent.EventAppeared.Event)
	var event events.CategoryAddedToMenu
	err := json.Unmarshal(recordedEvent.Data, &event)
	if err != nil {
		return err
	}
	err = menuEventHandler.menuRepository.AddCategoryToMenu(event.GetEntityID(), event.CategoryID)
	return err
}

func (menuEventHandler MenuEventHandler) HandleCategoryNameChanged(rawEvent *esdb.SubscriptionEvent) error {
	recordedEvent := eventutils2.DeserializeRecordedEvent(rawEvent.EventAppeared.Event)
	var event events.CategoryNameChanged
	err := json.Unmarshal(recordedEvent.Data, &event)
	if err != nil {
		return err
	}
	err = menuEventHandler.menuRepository.ChangeCategoryName(event.GetEntityID(), event.NewName)
	return err
}

func (menuEventHandler MenuEventHandler) HandleSubCategoryCreated(rawEvent *esdb.SubscriptionEvent) error {
	recordedEvent := eventutils2.DeserializeRecordedEvent(rawEvent.EventAppeared.Event)
	var event events.SubCategoryCreated
	err := json.Unmarshal(recordedEvent.Data, &event)
	if err != nil {
		return err
	}
	err = menuEventHandler.menuRepository.CreateSubCategory(event.GetEntityID(), event.Name)
	return err
}

func (menuEventHandler MenuEventHandler) HandleSubCategoryAddedToCategory(rawEvent *esdb.SubscriptionEvent) error {
	recordedEvent := eventutils2.DeserializeRecordedEvent(rawEvent.EventAppeared.Event)
	var event events.SubCategoryAddedToCategory
	err := json.Unmarshal(recordedEvent.Data, &event)
	if err != nil {
		return err
	}
	err = menuEventHandler.menuRepository.AddSubCategoryToCategory(event.GetEntityID(), event.SubCategoryID)
	return err
}

func (menuEventHandler MenuEventHandler) HandleMenuItemCreated(rawEvent *esdb.SubscriptionEvent) error {
	recordedEvent := eventutils2.DeserializeRecordedEvent(rawEvent.EventAppeared.Event)
	var event events.MenuItemCreated
	err := json.Unmarshal(recordedEvent.Data, &event)
	if err != nil {
		return err
	}
	err = menuEventHandler.menuRepository.CreateMenuItem(event.GetEntityID(), event.Name)
	return err
}

func (menuEventHandler MenuEventHandler) HandleMenuItemAddedToSubCategory(rawEvent *esdb.SubscriptionEvent) error {
	recordedEvent := eventutils2.DeserializeRecordedEvent(rawEvent.EventAppeared.Event)
	var event events.MenuItemAddedToSubCategory
	err := json.Unmarshal(recordedEvent.Data, &event)
	if err != nil {
		return err
	}
	err = menuEventHandler.menuRepository.AddMenuItemToSubCategory(event.GetEntityID(), event.MenuItemID)
	return err
}
