package entities

import (
	"encoding/json"

	"github.com/Resta-Inc/resta/pkg/events"
	"github.com/Resta-Inc/resta/pkg/eventutils"
	"github.com/Resta-Inc/resta/pkg/resources"
	"github.com/Resta-Inc/resta/pkg/utils"
	"github.com/gofrs/uuid"
)

// Models
type Menu struct {
	eventutils.Entity
	State MenuState
}

type MenuState struct {
	Name          string
	IsEnabled     bool
	CategoriesIDs []uuid.UUID
}

// Business Logic
func NewMenu() *Menu {
	menuID := utils.GenerateNewUUID()

	event := events.MenuCreated{
		EventInfo: eventutils.NewEventInfo(menuID),
		Name:      resources.DefaultMenuName("en"),
	}

	menu := &Menu{}
	menu.SetNew()
	eventutils.AddEvent(event, menu)
	return menu
}

func (menu Menu) GetName() string {
	return menu.State.Name
}

func (menu Menu) IsEnabled() bool {
	return menu.State.IsEnabled
}

func (menu Menu) GetCategoriesIDs() []uuid.UUID {
	return menu.State.CategoriesIDs
}

func (menu *Menu) Enable() {
	event := events.MenuEnabled{
		EventInfo: eventutils.NewEventInfo(menu.ID),
	}
	eventutils.AddEvent(event, menu)
}

func (menu *Menu) Disable() {
	event := events.MenuDisabled{
		EventInfo: eventutils.NewEventInfo(menu.ID),
	}
	eventutils.AddEvent(event, menu)
}

func (menu *Menu) ChangeName(newName string) {
	event := events.MenuNameChanged{
		EventInfo: eventutils.NewEventInfo(menu.ID),
		NewName:   newName,
	}
	eventutils.AddEvent(event, menu)
}

func (menu *Menu) AddCategory(categoryID uuid.UUID) {
	event := events.CategoryAddedToMenu{
		EventInfo:  eventutils.NewEventInfo(menu.ID),
		CategoryID: categoryID,
	}
	eventutils.AddEvent(event, menu)
}

// Events

func (menu *Menu) AppendEvent(event eventutils.IEvent) {
	menu.Events = append(menu.Events, event)
}

func (menu Menu) DeserializeEvent(event eventutils.Event) eventutils.IEvent {
	switch event.Name {
	case "MenuCreated":
		var e events.MenuCreated
		json.Unmarshal(event.Data, &e)
		return e
	case "MenuEnabled":
		var e events.MenuEnabled
		json.Unmarshal(event.Data, &e)
		return e
	case "MenuDisabled":
		var e events.MenuDisabled
		json.Unmarshal(event.Data, &e)
		return e
	case "MenuNameChanged":
		var e events.MenuNameChanged
		json.Unmarshal(event.Data, &e)
		return e
	case "CategoryAddedToMenu":
		var e events.CategoryAddedToMenu
		json.Unmarshal(event.Data, &e)
		return e
	}
	return nil
}

func (menu *Menu) ApplyEvent(event eventutils.IEvent) {
	eventType := utils.GetType(event)
	switch eventType {
	case "MenuCreated":
		applyMenuCreated(menu, event.(events.MenuCreated))
	case "MenuEnabled":
		applyMenuEnabled(menu)
	case "MenuDisabled":
		applyMenuDisabled(menu)
	case "MenuNameChanged":
		applyMenuNameChanged(menu, event.(events.MenuNameChanged))
	case "CategoryAddedToMenu":
		applyCategoryAddedToMenu(menu, event.(events.CategoryAddedToMenu))
	}
}

func applyMenuCreated(menu *Menu, e events.MenuCreated) {
	menu.State.Name = e.Name
	menu.ID = e.EntityID
}

func applyMenuEnabled(menu *Menu) {
	menu.State.IsEnabled = true
}

func applyMenuDisabled(menu *Menu) {
	menu.State.IsEnabled = false
}

func applyMenuNameChanged(menu *Menu, event events.MenuNameChanged) {
	menu.State.Name = event.NewName
}

func applyCategoryAddedToMenu(menu *Menu, event events.CategoryAddedToMenu) {
	menu.State.CategoriesIDs = append(menu.State.CategoriesIDs, event.CategoryID)
}
