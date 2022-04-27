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
	Name      string
	IsEnabled bool
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

func (menu Menu) IsEnabled() bool {
	return menu.State.IsEnabled
}

func (menu Menu) Enable() Menu {
	event := events.MenuEnabled{
		EntityEventInfo: eventutils.NewEntityEventInfo(menu.GetID()),
	}
	menu = eventutils.AddNewEvent(menu, event).(Menu)
	return menu
}

func (menu Menu) Disable() Menu {
	event := events.MenuDisabled{
		EntityEventInfo: eventutils.NewEntityEventInfo(menu.GetID()),
	}
	menu = eventutils.AddNewEvent(menu, event).(Menu)
	return menu
}

// Events
func (menu Menu) ApplyEvent(event eventutils.IEntityEvent) eventutils.IEntity {
	eventType := utils.GetType(event)
	switch eventType {
	case "MenuCreated":
		menu = menu.applyMenuCreated(event.(events.MenuCreated))
	case "MenuEnabled":
		menu = menu.applyMenuEnabled(event.(events.MenuEnabled))
	case "MenuDisabled":
		menu = menu.applyMenuDisabled(event.(events.MenuDisabled))
	}
	return menu
}

func (menu Menu) applyMenuCreated(event events.MenuCreated) Menu {
	menu.State.Name = event.Name
	menu.ID = event.EntityID
	return menu
}

func (menu Menu) applyMenuEnabled(event events.MenuEnabled) Menu {
	menu.State.IsEnabled = true
	return menu
}

func (menu Menu) applyMenuDisabled(event events.MenuDisabled) Menu {
	menu.State.IsEnabled = false
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
