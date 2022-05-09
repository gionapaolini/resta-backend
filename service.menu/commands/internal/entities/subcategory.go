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
	*eventutils.Entity
	State SubCategoryState
}
type SubCategoryState struct {
	Name string
}

// Business Logic
func NewSubCategory(categoryID uuid.UUID) SubCategory {
	subCategoryID := utils.GenerateNewUUID()

	event := events.SubCategoryCreated{
		EntityEventInfo:  eventutils.NewEntityEventInfo(subCategoryID),
		Name:             resources.DefaultSubCategoryName("en"),
		ParentCategoryID: categoryID,
	}

	subCategory := EmptySubCategory()
	return eventutils.AddNewEvent(subCategory, event).(SubCategory)
}

func EmptySubCategory() SubCategory {
	return SubCategory{
		Entity: &eventutils.Entity{},
	}
}

func (subCategory SubCategory) GetName() string {
	return subCategory.State.Name
}

func (subCategory SubCategory) ChangeName(newName string) SubCategory {
	event := events.SubCategoryNameChanged{
		EntityEventInfo: eventutils.NewEntityEventInfo(subCategory.ID),
		NewName:         newName,
	}
	return eventutils.AddNewEvent(subCategory, event).(SubCategory)
}

// Events
func (subCategory SubCategory) ApplyEvent(event eventutils.IEntityEvent) eventutils.IEntity {
	eventType := utils.GetType(event)
	switch eventType {
	case "SubCategoryCreated":
		subCategory = subCategory.applySubCategoryCreated(event.(events.SubCategoryCreated))
	case "SubCategoryNameChanged":
		subCategory = subCategory.applySubCategoryNameChanged(event.(events.SubCategoryNameChanged))
	}
	return subCategory
}

func (subCategory SubCategory) applySubCategoryCreated(event events.SubCategoryCreated) SubCategory {
	subCategory.ID = event.EntityID
	subCategory.State.Name = event.Name
	return subCategory
}

func (subCategory SubCategory) applySubCategoryNameChanged(event events.SubCategoryNameChanged) SubCategory {
	subCategory.State.Name = event.NewName
	return subCategory
}

func (subCategory SubCategory) DeserializeEvent(jsonData []byte) eventutils.IEvent {
	eventType, rawData := eventutils.GetRawDataFromSerializedEvent(jsonData)
	switch eventType {
	case "SubCategoryCreated":
		var e events.SubCategoryCreated
		json.Unmarshal(rawData, &e)
		return e
	case "SubCategoryNameChanged":
		var e events.SubCategoryNameChanged
		json.Unmarshal(rawData, &e)
		return e
	default:
		return nil
	}
}
