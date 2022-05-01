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
	Name     string
	ImageURL string
}

// Business Logic
func NewCategory(menuID uuid.UUID) Category {
	categoryID := utils.GenerateNewUUID()

	event := events.CategoryCreated{
		EntityEventInfo: eventutils.NewEntityEventInfo(categoryID),
		Name:            resources.DefaultCategoryName("en"),
		ImageURL:        resources.DefaultCategoryImageUrl(),
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

func (category Category) GetImageURL() string {
	return category.State.ImageURL
}

// Events
func (category Category) ApplyEvent(event eventutils.IEntityEvent) eventutils.IEntity {
	eventType := utils.GetType(event)
	switch eventType {
	case "CategoryCreated":
		category = category.applyCategoryCreated(event.(events.CategoryCreated))
	}
	return category
}

func (category Category) applyCategoryCreated(event events.CategoryCreated) Category {
	category.ID = event.EntityID
	category.State.Name = event.Name
	category.State.ImageURL = event.ImageURL
	return category
}

func (category Category) DeserializeEvent(jsonData []byte) eventutils.IEvent {
	eventType, rawData := eventutils.GetRawDataFromSerializedEvent(jsonData)
	switch eventType {
	case "CategoryCreated":
		var e events.CategoryCreated
		json.Unmarshal(rawData, &e)
		return e
	default:
		return nil
	}
}
