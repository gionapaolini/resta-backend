package entities

import (
	"encoding/json"

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

	event := MenuCreated{
		EntityEventInfo: eventutils.NewEntityEventInfo(menuID),
		Name:            resources.DefaultMenuName("en"),
	}

	menu := Menu{
		Entity: &eventutils.Entity{},
	}
	return eventutils.AddNewEvent(menu, event).(Menu)
}

func (menu Menu) GetName() string {
	return menu.State.Name
}

// Events
type MenuCreated struct {
	eventutils.EntityEventInfo
	Name string
}

func (event MenuCreated) Apply(entity eventutils.IEntity) eventutils.IEntity {
	menu := entity.(Menu)
	menu.State.Name = event.Name
	menu.ID = event.EntityID
	return menu
}

// Entity Logic
func (menu Menu) DeserializeEvent(jsonData []byte) eventutils.IEvent {
	eventType, rawData := eventutils.GetRawDataFromSerializedEvent(jsonData)
	switch eventType {
	case "MenuCreated":
		var e MenuCreated
		json.Unmarshal(rawData, &e)
		return e
	default:
		return nil
	}
}
