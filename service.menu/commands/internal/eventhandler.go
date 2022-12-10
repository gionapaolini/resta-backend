package internal

import (
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/Resta-Inc/resta/menu/commands/internal/entities"
	"github.com/Resta-Inc/resta/pkg/events"
	"github.com/Resta-Inc/resta/pkg/eventutils2"
)

type MenuEventHandler struct {
	entityRepository eventutils2.IEntityRepository
}

func NewMenuEventHandler(repo eventutils2.IEntityRepository) MenuEventHandler {
	return MenuEventHandler{
		entityRepository: repo,
	}
}

func (eventHandler MenuEventHandler) HandleCategoryCreated(rawEvent *esdb.SubscriptionEvent) error {
	event := eventutils2.DeserializeRecordedEvent(rawEvent.EventAppeared.Event)
	categoryCreatedEvent := entities.Category{}.DeserializeEvent(event).(events.CategoryCreated)
	menu, err := eventHandler.entityRepository.GetEntity(&entities.Menu{}, categoryCreatedEvent.ParentMenuID)
	if err != nil {
		return err
	}
	menu.(*entities.Menu).AddCategory(categoryCreatedEvent.GetEntityID())
	err = eventHandler.entityRepository.SaveEntity(menu)
	return err
}

func (eventHandler MenuEventHandler) HandleSubCategoryCreated(rawEvent *esdb.SubscriptionEvent) error {
	event := eventutils2.DeserializeRecordedEvent(rawEvent.EventAppeared.Event)
	subCategoryCreatedEvent := entities.SubCategory{}.DeserializeEvent(event).(events.SubCategoryCreated)
	category, err := eventHandler.entityRepository.GetEntity(&entities.Category{}, subCategoryCreatedEvent.ParentCategoryID)
	if err != nil {
		return err
	}
	category.(*entities.Category).AddSubCategory(subCategoryCreatedEvent.GetEntityID())
	err = eventHandler.entityRepository.SaveEntity(category)
	return err
}

func (eventHandler MenuEventHandler) HandleMenuItemCreated(rawEvent *esdb.SubscriptionEvent) error {
	event := eventutils2.DeserializeRecordedEvent(rawEvent.EventAppeared.Event)
	menuItemCreatedEvent := entities.MenuItem{}.DeserializeEvent(event).(events.MenuItemCreated)
	subCategory, err := eventHandler.entityRepository.GetEntity(&entities.SubCategory{}, menuItemCreatedEvent.ParentSubCategoryID)
	if err != nil {
		return err
	}

	subCategory.(*entities.SubCategory).AddMenuItem(menuItemCreatedEvent.GetEntityID())
	err = eventHandler.entityRepository.SaveEntity(subCategory)
	return err
}
