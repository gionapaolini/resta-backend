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
	*eventutils.Entity
	State CategoryState
}
type CategoryState struct {
	Name             string
	SubCategoriesIDs []uuid.UUID
}

// Business Logic
func NewCategory(menuID uuid.UUID) Category {
	categoryID := utils.GenerateNewUUID()

	event := events.CategoryCreated{
		EntityEventInfo: eventutils.NewEntityEventInfo(categoryID),
		Name:            resources.DefaultCategoryName("en"),
		ParentMenuID:    menuID,
	}

	category := EmptyCategory()
	return eventutils.AddNewEvent(category, event).(Category)
}

func EmptyCategory() Category {
	return Category{
		Entity: &eventutils.Entity{},
	}
}

func (category Category) GetName() string {
	return category.State.Name
}

func (category Category) GetSubCategoriesIDs() []uuid.UUID {
	return category.State.SubCategoriesIDs
}

func (category Category) ChangeName(newName string) Category {
	event := events.CategoryNameChanged{
		EntityEventInfo: eventutils.NewEntityEventInfo(category.ID),
		NewName:         newName,
	}
	return eventutils.AddNewEvent(category, event).(Category)
}

func (category Category) AddSubCategory(categoryID uuid.UUID) Category {
	event := events.SubCategoryAddedToCategory{
		EntityEventInfo: eventutils.NewEntityEventInfo(category.GetID()),
		SubCategoryID:   categoryID,
	}
	category = eventutils.AddNewEvent(category, event).(Category)
	return category
}

// Events
func (category Category) ApplyEvent(event eventutils.IEntityEvent) eventutils.IEntity {
	eventType := utils.GetType(event)
	switch eventType {
	case "CategoryCreated":
		category = category.applyCategoryCreated(event.(events.CategoryCreated))
	case "CategoryNameChanged":
		category = category.applyCategoryNameChanged(event.(events.CategoryNameChanged))
	case "SubCategoryAddedToCategory":
		category = category.applySubCategoryAddedToCategory(event.(events.SubCategoryAddedToCategory))
	}
	return category
}

func (category Category) applyCategoryCreated(event events.CategoryCreated) Category {
	category.ID = event.EntityID
	category.State.Name = event.Name
	return category
}

func (category Category) applyCategoryNameChanged(event events.CategoryNameChanged) Category {
	category.State.Name = event.NewName
	return category
}

func (category Category) applySubCategoryAddedToCategory(event events.SubCategoryAddedToCategory) Category {
	category.State.SubCategoriesIDs = append(category.State.SubCategoriesIDs, event.SubCategoryID)
	return category
}

func (category Category) DeserializeEvent(jsonData []byte) eventutils.IEvent {
	eventType, rawData := eventutils.GetRawDataFromSerializedEvent(jsonData)
	switch eventType {
	case "CategoryCreated":
		var e events.CategoryCreated
		json.Unmarshal(rawData, &e)
		return e
	case "CategoryNameChanged":
		var e events.CategoryNameChanged
		json.Unmarshal(rawData, &e)
		return e
	case "SubCategoryAddedToCategory":
		var e events.SubCategoryAddedToCategory
		json.Unmarshal(rawData, &e)
		return e
	default:
		return nil
	}
}
