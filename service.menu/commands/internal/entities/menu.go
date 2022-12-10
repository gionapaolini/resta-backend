package entities

import (
	"encoding/json"

	"github.com/Resta-Inc/resta/pkg/events2"
	"github.com/Resta-Inc/resta/pkg/eventutils2"
	"github.com/Resta-Inc/resta/pkg/resources"
	"github.com/Resta-Inc/resta/pkg/utils"
	"github.com/gofrs/uuid"
)

// Models
type Menu struct {
	eventutils2.Entity
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

	event := events2.MenuCreated{
		EventInfo: eventutils2.NewEventInfo(menuID),
		Name:      resources.DefaultMenuName("en"),
	}

	menu := &Menu{}
	menu.SetNew()
	eventutils2.AddEvent(event, menu)
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
	event := events2.MenuEnabled{
		EventInfo: eventutils2.NewEventInfo(menu.ID),
	}
	eventutils2.AddEvent(event, menu)
}

func (menu *Menu) Disable() {
	event := events2.MenuDisabled{
		EventInfo: eventutils2.NewEventInfo(menu.ID),
	}
	eventutils2.AddEvent(event, menu)
}

func (menu *Menu) ChangeName(newName string) {
	event := events2.MenuNameChanged{
		EventInfo: eventutils2.NewEventInfo(menu.ID),
		NewName:   newName,
	}
	eventutils2.AddEvent(event, menu)
}

func (menu *Menu) AddCategory(categoryID uuid.UUID) {
	event := events2.CategoryAddedToMenu{
		EventInfo:  eventutils2.NewEventInfo(menu.ID),
		CategoryID: categoryID,
	}
	eventutils2.AddEvent(event, menu)
}

// Events

func (menu *Menu) AppendEvent(event eventutils2.IEvent) {
	menu.Events = append(menu.Events, event)
}

func (menu Menu) DeserializeEvent(event eventutils2.Event) eventutils2.IEvent {
	switch event.Name {
	case "MenuCreated":
		var e events2.MenuCreated
		json.Unmarshal(event.Data, &e)
		return e
	case "MenuEnabled":
		var e events2.MenuEnabled
		json.Unmarshal(event.Data, &e)
		return e
	case "MenuDisabled":
		var e events2.MenuDisabled
		json.Unmarshal(event.Data, &e)
		return e
	case "MenuNameChanged":
		var e events2.MenuNameChanged
		json.Unmarshal(event.Data, &e)
		return e
	case "CategoryAddedToMenu":
		var e events2.CategoryAddedToMenu
		json.Unmarshal(event.Data, &e)
		return e
	}
	return nil
}

func (menu *Menu) ApplyEvent(event eventutils2.IEvent) {
	eventType := utils.GetType(event)
	switch eventType {
	case "MenuCreated":
		applyMenuCreated(menu, event.(events2.MenuCreated))
	case "MenuEnabled":
		applyMenuEnabled(menu)
	case "MenuDisabled":
		applyMenuDisabled(menu)
	case "MenuNameChanged":
		applyMenuNameChanged(menu, event.(events2.MenuNameChanged))
	case "CategoryAddedToMenu":
		applyCategoryAddedToMenu(menu, event.(events2.CategoryAddedToMenu))
	}
}

func applyMenuCreated(menu *Menu, e events2.MenuCreated) {
	menu.State.Name = e.Name
	menu.ID = e.EntityID
}

func applyMenuEnabled(menu *Menu) {
	menu.State.IsEnabled = true
}

func applyMenuDisabled(menu *Menu) {
	menu.State.IsEnabled = false
}

func applyMenuNameChanged(menu *Menu, event events2.MenuNameChanged) {
	menu.State.Name = event.NewName
}

func applyCategoryAddedToMenu(menu *Menu, event events2.CategoryAddedToMenu) {
	menu.State.CategoriesIDs = append(menu.State.CategoriesIDs, event.CategoryID)
}
