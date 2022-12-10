package entities

import (
	"encoding/json"
	"time"

	"github.com/Resta-Inc/resta/pkg/events"
	"github.com/Resta-Inc/resta/pkg/eventutils"
	"github.com/Resta-Inc/resta/pkg/resources"
	"github.com/Resta-Inc/resta/pkg/utils"
	"github.com/gofrs/uuid"
)

type MenuItem struct {
	eventutils.Entity
	State MenuItemState
}
type MenuItemState struct {
	Name                     string
	EstimatedPreparationTime time.Duration
}

// Business Logic
func NewMenuItem(subCategoryID uuid.UUID) *MenuItem {
	categoryID := utils.GenerateNewUUID()

	event := events.MenuItemCreated{
		EventInfo:           eventutils.NewEventInfo(categoryID),
		Name:                resources.DefaultMenuItemName("en"),
		ParentSubCategoryID: subCategoryID,
	}

	menuItem := &MenuItem{}
	menuItem.SetNew()
	eventutils.AddEvent(event, menuItem)
	return menuItem
}

func (menuItem MenuItem) GetName() string {
	return menuItem.State.Name
}

func (menuItem MenuItem) GetEstimatedPreparationtime() time.Duration {
	return menuItem.State.EstimatedPreparationTime
}

func (menuItem *MenuItem) ChangeName(newName string) {
	event := events.MenuItemNameChanged{
		EventInfo: eventutils.NewEventInfo(menuItem.GetID()),
		NewName:   newName,
	}
	eventutils.AddEvent(event, menuItem)
}

func (menuItem *MenuItem) ChangeEstimatedPreparationTime(newTime time.Duration) {
	event := events.MenuItemEstimatedPreparationTimeChanged{
		EventInfo:   eventutils.NewEventInfo(menuItem.GetID()),
		NewEstimate: newTime,
	}

	eventutils.AddEvent(event, menuItem)
}

// Events

func (menuItem *MenuItem) AppendEvent(event eventutils.IEvent) {
	menuItem.Events = append(menuItem.Events, event)
}

func (menuItem MenuItem) DeserializeEvent(event eventutils.Event) eventutils.IEvent {
	switch event.Name {
	case "MenuItemCreated":
		var e events.MenuItemCreated
		json.Unmarshal(event.Data, &e)
		return e
	case "MenuItemNameChanged":
		var e events.MenuItemNameChanged
		json.Unmarshal(event.Data, &e)
		return e
	case "MenuItemEstimatedPreparationTimeChanged":
		var e events.MenuItemEstimatedPreparationTimeChanged
		json.Unmarshal(event.Data, &e)
		return e
	default:
		return nil
	}
}

func (menuItem *MenuItem) ApplyEvent(event eventutils.IEvent) {
	eventType := utils.GetType(event)
	switch eventType {
	case "MenuItemCreated":
		applyMenuItemCreated(menuItem, event.(events.MenuItemCreated))
	case "MenuItemNameChanged":
		applyMenuItemNameChanged(menuItem, event.(events.MenuItemNameChanged))
	case "MenuItemEstimatedPreparationTimeChanged":
		applyMenuItemEstimatedPreparationTimeChanged(menuItem, event.(events.MenuItemEstimatedPreparationTimeChanged))

	}
}

func applyMenuItemCreated(menuItem *MenuItem, event events.MenuItemCreated) {
	menuItem.ID = event.EntityID
	menuItem.State.Name = event.Name
}

func applyMenuItemNameChanged(menuItem *MenuItem, event events.MenuItemNameChanged) {
	menuItem.ID = event.EntityID
	menuItem.State.Name = event.NewName
}

func applyMenuItemEstimatedPreparationTimeChanged(menuItem *MenuItem, event events.MenuItemEstimatedPreparationTimeChanged) {
	menuItem.State.EstimatedPreparationTime = event.NewEstimate
}
