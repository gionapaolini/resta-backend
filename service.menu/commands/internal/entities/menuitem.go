package entities

import (
	"encoding/json"

	"github.com/Resta-Inc/resta/pkg/events"
	"github.com/Resta-Inc/resta/pkg/eventutils"
	"github.com/Resta-Inc/resta/pkg/resources"
	"github.com/Resta-Inc/resta/pkg/utils"
	"github.com/gofrs/uuid"
)

type MenuItem struct {
	*eventutils.Entity
	State MenuItemState
}
type MenuItemState struct {
	Name             string
	SubCategoriesIDs []uuid.UUID
}

// Business Logic
func NewMenuItem(subCategoryID uuid.UUID) MenuItem {
	categoryID := utils.GenerateNewUUID()

	event := events.MenuItemCreated{
		EntityEventInfo:     eventutils.NewEntityEventInfo(categoryID),
		Name:                resources.DefaultMenuItemName("en"),
		ParentSubCategoryID: subCategoryID,
	}

	category := EmptyMenuItem()
	return eventutils.AddNewEvent(category, event).(MenuItem)
}

func EmptyMenuItem() MenuItem {
	return MenuItem{
		Entity: &eventutils.Entity{},
	}
}

func (menuItem MenuItem) GetName() string {
	return menuItem.State.Name
}

func (menuItem MenuItem) ChangeName(newName string) MenuItem {
	event := events.MenuItemNameChanged{
		EntityEventInfo: eventutils.NewEntityEventInfo(menuItem.GetID()),
		NewName:         newName,
	}

	return eventutils.AddNewEvent(menuItem, event).(MenuItem)
}

// Events
func (menuItem MenuItem) ApplyEvent(event eventutils.IEntityEvent) eventutils.IEntity {
	eventType := utils.GetType(event)
	switch eventType {
	case "MenuItemCreated":
		menuItem = menuItem.applyMenuItemCreated(event.(events.MenuItemCreated))
	case "MenuItemNameChanged":
		menuItem = menuItem.applyMenuItemNameChanged(event.(events.MenuItemNameChanged))

	}
	return menuItem
}

func (menuItem MenuItem) applyMenuItemCreated(event events.MenuItemCreated) MenuItem {
	menuItem.ID = event.EntityID
	menuItem.State.Name = event.Name
	return menuItem
}

func (menuItem MenuItem) applyMenuItemNameChanged(event events.MenuItemNameChanged) MenuItem {
	menuItem.ID = event.EntityID
	menuItem.State.Name = event.NewName
	return menuItem
}

func (menuItem MenuItem) DeserializeEvent(jsonData []byte) eventutils.IEvent {
	eventType, rawData := eventutils.GetRawDataFromSerializedEvent(jsonData)
	switch eventType {
	case "MenuItemCreated":
		var e events.MenuItemCreated
		json.Unmarshal(rawData, &e)
		return e
	case "MenuItemNameChanged":
		var e events.MenuItemNameChanged
		json.Unmarshal(rawData, &e)
		return e
	default:
		return nil
	}
}
