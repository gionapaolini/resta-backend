package entities

import (
	"encoding/json"

	"github.com/Resta-Inc/resta/pkg/events"
	"github.com/Resta-Inc/resta/pkg/eventutils"
	"github.com/Resta-Inc/resta/pkg/resources"
	"github.com/Resta-Inc/resta/pkg/utils"
	"github.com/gofrs/uuid"
)

type Category struct {
	eventutils.Entity
	State CategoryState
}
type CategoryState struct {
	Name             string
	SubCategoriesIDs []uuid.UUID
}

// Business Logic
func NewCategory(menuID uuid.UUID) *Category {
	categoryID := utils.GenerateNewUUID()

	event := events.CategoryCreated{
		EventInfo:    eventutils.NewEventInfo(categoryID),
		Name:         resources.DefaultCategoryName("en"),
		ParentMenuID: menuID,
	}

	category := &Category{}
	category.SetNew()
	eventutils.AddEvent(event, category)
	return category
}

func (category Category) GetName() string {
	return category.State.Name
}

func (category Category) GetSubCategoriesIDs() []uuid.UUID {
	return category.State.SubCategoriesIDs
}

func (category *Category) ChangeName(newName string) {
	event := events.CategoryNameChanged{
		EventInfo: eventutils.NewEventInfo(category.ID),
		NewName:   newName,
	}
	eventutils.AddEvent(event, category)
}

func (category *Category) AddSubCategory(categoryID uuid.UUID) {
	event := events.SubCategoryAddedToCategory{
		EventInfo:     eventutils.NewEventInfo(category.GetID()),
		SubCategoryID: categoryID,
	}
	eventutils.AddEvent(event, category)
}

// Events
func (category *Category) AppendEvent(event eventutils.IEvent) {
	category.Events = append(category.Events, event)
}

func (category Category) DeserializeEvent(event eventutils.Event) eventutils.IEvent {
	switch event.Name {
	case "CategoryCreated":
		var e events.CategoryCreated
		json.Unmarshal(event.Data, &e)
		return e
	case "CategoryNameChanged":
		var e events.CategoryNameChanged
		json.Unmarshal(event.Data, &e)
		return e
	case "SubCategoryAddedToCategory":
		var e events.SubCategoryAddedToCategory
		json.Unmarshal(event.Data, &e)
		return e
	default:
		return nil
	}
}

func (category *Category) ApplyEvent(event eventutils.IEvent) {
	eventType := utils.GetType(event)
	switch eventType {
	case "CategoryCreated":
		applyCategoryCreated(category, event.(events.CategoryCreated))
	case "CategoryNameChanged":
		applyCategoryNameChanged(category, event.(events.CategoryNameChanged))
	case "SubCategoryAddedToCategory":
		applySubCategoryAddedToCategory(category, event.(events.SubCategoryAddedToCategory))
	}
}

func applyCategoryCreated(category *Category, event events.CategoryCreated) {
	category.ID = event.EntityID
	category.State.Name = event.Name
}

func applyCategoryNameChanged(category *Category, event events.CategoryNameChanged) {
	category.State.Name = event.NewName
}

func applySubCategoryAddedToCategory(category *Category, event events.SubCategoryAddedToCategory) {
	category.State.SubCategoriesIDs = append(category.State.SubCategoriesIDs, event.SubCategoryID)
}
