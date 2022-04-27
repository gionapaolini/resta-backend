package entities

import (
	"encoding/json"

	"github.com/Resta-Inc/resta/pkg/events"
	"github.com/Resta-Inc/resta/pkg/eventutils"
	"github.com/Resta-Inc/resta/pkg/resources"
	"github.com/Resta-Inc/resta/pkg/utils"
)

// Models
type Menu struct {
	*eventutils.Entity
	State MenuState
}

type MenuState struct {
	Name string
}

// Business Logic
func NewMenu() Menu {
	menuID := utils.GenerateNewUUID()

	event := events.MenuCreated{
		EntityEventInfo: eventutils.NewEntityEventInfo(menuID),
		Name:            resources.DefaultMenuName("en"),
	}

	menu := EmptyMenu()
	return eventutils.AddNewEvent(menu, event).(Menu)
}

func EmptyMenu() Menu {
	return Menu{
		Entity: &eventutils.Entity{},
	}
}

func (menu Menu) GetName() string {
	return menu.State.Name
}

// Events
func (menu Menu) ApplyEvent(event eventutils.IEntityEvent) eventutils.IEntity {
	eventType := utils.GetType(event)
	switch eventType {
	case "MenuCreated":
		menu = menu.applyMenuCreated(event.(events.MenuCreated))
	}
	return menu
}

func (menu Menu) applyMenuCreated(event events.MenuCreated) Menu {
	menu.State.Name = event.Name
	menu.ID = event.EntityID
	return menu
}

func (menu Menu) DeserializeEvent(jsonData []byte) eventutils.IEvent {
	eventType, rawData := eventutils.GetRawDataFromSerializedEvent(jsonData)
	switch eventType {
	case "MenuCreated":
		var e events.MenuCreated
		json.Unmarshal(rawData, &e)
		return e
	default:
		return nil
	}
}
