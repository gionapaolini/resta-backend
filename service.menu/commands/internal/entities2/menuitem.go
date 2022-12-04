package entities2

import (
	"encoding/json"
	"time"

	"github.com/Resta-Inc/resta/pkg/events2"
	"github.com/Resta-Inc/resta/pkg/eventutils2"
	"github.com/Resta-Inc/resta/pkg/resources"
	"github.com/Resta-Inc/resta/pkg/utils"
	"github.com/gofrs/uuid"
)

type MenuItem struct {
	eventutils2.Entity
	State MenuItemState
}
type MenuItemState struct {
	Name                     string
	EstimatedPreparationTime time.Duration
}

// Business Logic
func NewMenuItem(subCategoryID uuid.UUID) *MenuItem {
	categoryID := utils.GenerateNewUUID()

	event := events2.MenuItemCreated{
		EventInfo:           eventutils2.NewEventInfo(categoryID),
		Name:                resources.DefaultMenuItemName("en"),
		ParentSubCategoryID: subCategoryID,
	}

	menuItem := &MenuItem{}
	eventutils2.AddEvent(event, menuItem)
	return menuItem
}

func (menuItem MenuItem) GetName() string {
	return menuItem.State.Name
}

func (menuItem MenuItem) GetEstimatedPreparationtime() time.Duration {
	return menuItem.State.EstimatedPreparationTime
}

func (menuItem *MenuItem) ChangeName(newName string) {
	event := events2.MenuItemNameChanged{
		EventInfo: eventutils2.NewEventInfo(menuItem.GetID()),
		NewName:   newName,
	}
	eventutils2.AddEvent(event, menuItem)
}

func (menuItem *MenuItem) ChangeEstimatedPreparationTime(newTime time.Duration) {
	event := events2.MenuItemEstimatedPreparationTimeChanged{
		EventInfo:   eventutils2.NewEventInfo(menuItem.GetID()),
		NewEstimate: newTime,
	}

	eventutils2.AddEvent(event, menuItem)
}

// Events

func (menuItem *MenuItem) AppendEvent(event eventutils2.IEvent) {
	menuItem.Events = append(menuItem.Events, event)
}

func (menuItem MenuItem) DeserializeEvent(event eventutils2.Event) eventutils2.IEvent {
	switch event.Name {
	case "MenuItemCreated":
		var e events2.MenuItemCreated
		json.Unmarshal(event.Data, &e)
		return e
	case "MenuItemNameChanged":
		var e events2.MenuItemNameChanged
		json.Unmarshal(event.Data, &e)
		return e
	case "MenuItemEstimatedPreparationTimeChanged":
		var e events2.MenuItemEstimatedPreparationTimeChanged
		json.Unmarshal(event.Data, &e)
		return e
	default:
		return nil
	}
}

func (menuItem *MenuItem) ApplyEvent(event eventutils2.IEvent) {
	eventType := utils.GetType(event)
	switch eventType {
	case "MenuItemCreated":
		applyMenuItemCreated(menuItem, event.(events2.MenuItemCreated))
	case "MenuItemNameChanged":
		applyMenuItemNameChanged(menuItem, event.(events2.MenuItemNameChanged))
	case "MenuItemEstimatedPreparationTimeChanged":
		applyMenuItemEstimatedPreparationTimeChanged(menuItem, event.(events2.MenuItemEstimatedPreparationTimeChanged))

	}
}

func applyMenuItemCreated(menuItem *MenuItem, event events2.MenuItemCreated) {
	menuItem.ID = event.EntityID
	menuItem.State.Name = event.Name
}

func applyMenuItemNameChanged(menuItem *MenuItem, event events2.MenuItemNameChanged) {
	menuItem.ID = event.EntityID
	menuItem.State.Name = event.NewName
}

func applyMenuItemEstimatedPreparationTimeChanged(menuItem *MenuItem, event events2.MenuItemEstimatedPreparationTimeChanged) {
	menuItem.State.EstimatedPreparationTime = event.NewEstimate
}
