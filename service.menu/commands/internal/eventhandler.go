package internal

import (
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/Resta-Inc/resta/menu/commands/internal/entities2"
	"github.com/Resta-Inc/resta/pkg/events2"
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
	categoryCreatedEvent := entities2.Category{}.DeserializeEvent(event).(events2.CategoryCreated)
	menu, err := eventHandler.entityRepository.GetEntity(&entities2.Menu{}, categoryCreatedEvent.ParentMenuID)
	if err != nil {
		return err
	}
	menu.(*entities2.Menu).AddCategory(categoryCreatedEvent.GetEntityID())
	err = eventHandler.entityRepository.SaveEntity(menu)
	return err
}

func (eventHandler MenuEventHandler) HandleSubCategoryCreated(rawEvent *esdb.SubscriptionEvent) error {
	event := eventutils2.DeserializeRecordedEvent(rawEvent.EventAppeared.Event)
	subCategoryCreatedEvent := entities2.SubCategory{}.DeserializeEvent(event).(events2.SubCategoryCreated)
	category, err := eventHandler.entityRepository.GetEntity(&entities2.Category{}, subCategoryCreatedEvent.ParentCategoryID)
	if err != nil {
		return err
	}
	category.(*entities2.Category).AddSubCategory(subCategoryCreatedEvent.GetEntityID())
	err = eventHandler.entityRepository.SaveEntity(category)
	return err
}

func (eventHandler MenuEventHandler) HandleMenuItemCreated(rawEvent *esdb.SubscriptionEvent) error {
	event := eventutils2.DeserializeRecordedEvent(rawEvent.EventAppeared.Event)
	menuItemCreatedEvent := entities2.MenuItem{}.DeserializeEvent(event).(events2.MenuItemCreated)
	subCategory, err := eventHandler.entityRepository.GetEntity(&entities2.SubCategory{}, menuItemCreatedEvent.ParentSubCategoryID)
	if err != nil {
		return err
	}

	subCategory.(*entities2.SubCategory).AddMenuItem(menuItemCreatedEvent.GetEntityID())
	err = eventHandler.entityRepository.SaveEntity(subCategory)
	return err
}
