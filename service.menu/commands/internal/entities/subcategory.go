package entities

import (
	"encoding/json"

	"github.com/Resta-Inc/resta/pkg/events"
	"github.com/Resta-Inc/resta/pkg/eventutils"
	"github.com/Resta-Inc/resta/pkg/resources"
	"github.com/Resta-Inc/resta/pkg/utils"
	"github.com/gofrs/uuid"
)

type SubCategory struct {
	eventutils.Entity
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
		EventInfo:        eventutils.NewEventInfo(subCategoryID),
		Name:             resources.DefaultSubCategoryName("en"),
		ParentCategoryID: categoryID,
	}

	subCategory := &SubCategory{}
	subCategory.SetNew()
	eventutils.AddEvent(event, subCategory)
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
		EventInfo: eventutils.NewEventInfo(subCategory.ID),
		NewName:   newName,
	}
	eventutils.AddEvent(event, subCategory)
}

func (subCategory *SubCategory) AddMenuItem(menuItemID uuid.UUID) {
	event := events.MenuItemAddedToSubCategory{
		EventInfo:  eventutils.NewEventInfo(subCategory.ID),
		MenuItemID: menuItemID,
	}
	eventutils.AddEvent(event, subCategory)
}

// Events

func (subCategory SubCategory) DeserializeEvent(event eventutils.Event) eventutils.IEvent {
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

func (subCategory *SubCategory) ApplyEvent(event eventutils.IEvent) {
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
