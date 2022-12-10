package entities

import (
	"encoding/json"

	"github.com/Resta-Inc/resta/pkg/events2"
	"github.com/Resta-Inc/resta/pkg/eventutils2"
	"github.com/Resta-Inc/resta/pkg/resources"
	"github.com/Resta-Inc/resta/pkg/utils"
	"github.com/gofrs/uuid"
)

type Category struct {
	eventutils2.Entity
	State CategoryState
}
type CategoryState struct {
	Name             string
	SubCategoriesIDs []uuid.UUID
}

// Business Logic
func NewCategory(menuID uuid.UUID) *Category {
	categoryID := utils.GenerateNewUUID()

	event := events2.CategoryCreated{
		EventInfo:    eventutils2.NewEventInfo(categoryID),
		Name:         resources.DefaultCategoryName("en"),
		ParentMenuID: menuID,
	}

	category := &Category{}
	category.SetNew()
	eventutils2.AddEvent(event, category)
	return category
}

func (category Category) GetName() string {
	return category.State.Name
}

func (category Category) GetSubCategoriesIDs() []uuid.UUID {
	return category.State.SubCategoriesIDs
}

func (category *Category) ChangeName(newName string) {
	event := events2.CategoryNameChanged{
		EventInfo: eventutils2.NewEventInfo(category.ID),
		NewName:   newName,
	}
	eventutils2.AddEvent(event, category)
}

func (category *Category) AddSubCategory(categoryID uuid.UUID) {
	event := events2.SubCategoryAddedToCategory{
		EventInfo:     eventutils2.NewEventInfo(category.GetID()),
		SubCategoryID: categoryID,
	}
	eventutils2.AddEvent(event, category)
}

// Events
func (category *Category) AppendEvent(event eventutils2.IEvent) {
	category.Events = append(category.Events, event)
}

func (category Category) DeserializeEvent(event eventutils2.Event) eventutils2.IEvent {
	switch event.Name {
	case "CategoryCreated":
		var e events2.CategoryCreated
		json.Unmarshal(event.Data, &e)
		return e
	case "CategoryNameChanged":
		var e events2.CategoryNameChanged
		json.Unmarshal(event.Data, &e)
		return e
	case "SubCategoryAddedToCategory":
		var e events2.SubCategoryAddedToCategory
		json.Unmarshal(event.Data, &e)
		return e
	default:
		return nil
	}
}

func (category *Category) ApplyEvent(event eventutils2.IEvent) {
	eventType := utils.GetType(event)
	switch eventType {
	case "CategoryCreated":
		applyCategoryCreated(category, event.(events2.CategoryCreated))
	case "CategoryNameChanged":
		applyCategoryNameChanged(category, event.(events2.CategoryNameChanged))
	case "SubCategoryAddedToCategory":
		applySubCategoryAddedToCategory(category, event.(events2.SubCategoryAddedToCategory))
	}
}

func applyCategoryCreated(category *Category, event events2.CategoryCreated) {
	category.ID = event.EntityID
	category.State.Name = event.Name
}

func applyCategoryNameChanged(category *Category, event events2.CategoryNameChanged) {
	category.State.Name = event.NewName
}

func applySubCategoryAddedToCategory(category *Category, event events2.SubCategoryAddedToCategory) {
	category.State.SubCategoriesIDs = append(category.State.SubCategoriesIDs, event.SubCategoryID)
}
