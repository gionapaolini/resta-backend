package entities

import (
	"encoding/json"

	"github.com/Resta-Inc/resta/pkg/events"
	"github.com/Resta-Inc/resta/pkg/eventutils2"
	"github.com/Resta-Inc/resta/pkg/resources"
	"github.com/Resta-Inc/resta/pkg/utils"
	"github.com/gofrs/uuid"
)

type SubCategory struct {
	eventutils2.Entity
	State SubCategoryState
}
type SubCategoryState struct {
	Name         string
	MenuItemsIDs []uuid.UUID
}

// Business Logic
func NewSubCategory(categoryID uuid.UUID) *SubCategory {
	subCategoryID := utils.GenerateNewUUID()

	event := events.SubCategoryCreated{
		EventInfo:        eventutils2.NewEventInfo(subCategoryID),
		Name:             resources.DefaultSubCategoryName("en"),
		ParentCategoryID: categoryID,
	}

	subCategory := &SubCategory{}
	subCategory.SetNew()
	eventutils2.AddEvent(event, subCategory)
	return subCategory
}

func (subCategory SubCategory) GetName() string {
	return subCategory.State.Name
}

func (subCategory SubCategory) GetMenuItemsIDs() []uuid.UUID {
	return subCategory.State.MenuItemsIDs
}

func (subCategory *SubCategory) ChangeName(newName string) {
	event := events.SubCategoryNameChanged{
		EventInfo: eventutils2.NewEventInfo(subCategory.ID),
		NewName:   newName,
	}
	eventutils2.AddEvent(event, subCategory)
}

func (subCategory *SubCategory) AddMenuItem(menuItemID uuid.UUID) {
	event := events.MenuItemAddedToSubCategory{
		EventInfo:  eventutils2.NewEventInfo(subCategory.ID),
		MenuItemID: menuItemID,
	}
	eventutils2.AddEvent(event, subCategory)
}

// Events

func (subCategory *SubCategory) AppendEvent(event eventutils2.IEvent) {
	subCategory.Events = append(subCategory.Events, event)
}

func (subCategory SubCategory) DeserializeEvent(event eventutils2.Event) eventutils2.IEvent {
	switch event.Name {
	case "SubCategoryCreated":
		var e events.SubCategoryCreated
		json.Unmarshal(event.Data, &e)
		return e
	case "SubCategoryNameChanged":
		var e events.SubCategoryNameChanged
		json.Unmarshal(event.Data, &e)
		return e
	case "MenuItemAddedToSubCategory":
		var e events.MenuItemAddedToSubCategory
		json.Unmarshal(event.Data, &e)
		return e
	default:
		return nil
	}
}

func (subCategory *SubCategory) ApplyEvent(event eventutils2.IEvent) {
	eventType := utils.GetType(event)
	switch eventType {
	case "SubCategoryCreated":
		applySubCategoryCreated(subCategory, event.(events.SubCategoryCreated))
	case "SubCategoryNameChanged":
		applySubCategoryNameChanged(subCategory, event.(events.SubCategoryNameChanged))
	case "MenuItemAddedToSubCategory":
		applyMenuItemAddedToSubCategory(subCategory, event.(events.MenuItemAddedToSubCategory))
	}
}

func applySubCategoryCreated(subCategory *SubCategory, event events.SubCategoryCreated) {
	subCategory.ID = event.EntityID
	subCategory.State.Name = event.Name
}

func applySubCategoryNameChanged(subCategory *SubCategory, event events.SubCategoryNameChanged) {
	subCategory.State.Name = event.NewName
}

func applyMenuItemAddedToSubCategory(subCategory *SubCategory, event events.MenuItemAddedToSubCategory) {
	subCategory.State.MenuItemsIDs = append(subCategory.State.MenuItemsIDs, event.MenuItemID)
}
