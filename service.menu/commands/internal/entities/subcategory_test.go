package entities

import (
	"testing"

	"github.com/Resta-Inc/resta/pkg/events"
	"github.com/Resta-Inc/resta/pkg/eventutils"
	"github.com/Resta-Inc/resta/pkg/resources"
	"github.com/Resta-Inc/resta/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestCreateSubCategory(t *testing.T) {
	// Act
	categoryID := utils.GenerateNewUUID()
	subCategory := NewSubCategory(categoryID)

	// Assert
	latestEvent := subCategory.Events[len(subCategory.Events)-1]
	require.True(t, subCategory.IsNew())
	require.Equal(t, resources.DefaultSubCategoryName("en"), subCategory.GetName())
	require.Len(t, subCategory.Events, 1)
	require.IsType(t, events.SubCategoryCreated{}, latestEvent)
	require.Equal(t, utils.Time.Now(), latestEvent.GetTimeStamp())
	require.Equal(t, subCategory.ID, latestEvent.GetEntityID())
	require.False(t, subCategory.IsDeleted)
}

func TestChangeSubCategoryName(t *testing.T) {
	// Arrange
	categoryID := utils.GenerateNewUUID()
	subCategory := NewSubCategory(categoryID)
	newName := "New name"

	// Act
	subCategory.ChangeName(newName)

	// Assert
	latestEvent := subCategory.Events[len(subCategory.Events)-1]
	require.Equal(t, newName, subCategory.GetName())
	require.IsType(t, events.SubCategoryNameChanged{}, latestEvent)
}

func Test_AddMenuItem(t *testing.T) {
	// Arrange
	menuItemID := utils.GenerateNewUUID()
	subCategory := NewSubCategory(utils.GenerateNewUUID())

	// Act
	subCategory.AddMenuItem(menuItemID)

	// Assert
	latestEvent := subCategory.Events[len(subCategory.Events)-1]
	require.Contains(t, subCategory.GetMenuItemsIDs(), menuItemID)
	require.IsType(t, events.MenuItemAddedToSubCategory{}, latestEvent)
}

func Test_DeserializeSubCategoryEvent(t *testing.T) {
	// Arrange
	events := []eventutils.IEvent{
		events.SubCategoryCreated{
			EventInfo: eventutils.NewEventInfo(utils.GenerateNewUUID()),
		},
		events.SubCategoryNameChanged{
			EventInfo: eventutils.NewEventInfo(utils.GenerateNewUUID()),
		},
		events.MenuItemAddedToSubCategory{
			EventInfo: eventutils.NewEventInfo(utils.GenerateNewUUID()),
		},
	}

	for _, event := range events {
		serialized := eventutils.SerializedEvent(event)

		// Act
		deserialized := SubCategory{}.DeserializeEvent(serialized)

		// Assert
		require.Equal(t, event, deserialized)
	}
}
