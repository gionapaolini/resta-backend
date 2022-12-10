package entities

import (
	"testing"

	"github.com/Resta-Inc/resta/pkg/events2"
	"github.com/Resta-Inc/resta/pkg/eventutils2"
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
	require.IsType(t, events2.SubCategoryCreated{}, latestEvent)
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
	require.IsType(t, events2.SubCategoryNameChanged{}, latestEvent)
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
	require.IsType(t, events2.MenuItemAddedToSubCategory{}, latestEvent)
}

func Test_DeserializeSubCategoryEvent(t *testing.T) {
	// Arrange
	events := []eventutils2.IEvent{
		events2.SubCategoryCreated{
			EventInfo: eventutils2.NewEventInfo(utils.GenerateNewUUID()),
		},
		events2.SubCategoryNameChanged{
			EventInfo: eventutils2.NewEventInfo(utils.GenerateNewUUID()),
		},
		events2.MenuItemAddedToSubCategory{
			EventInfo: eventutils2.NewEventInfo(utils.GenerateNewUUID()),
		},
	}

	for _, event := range events {
		serialized := eventutils2.SerializedEvent(event)

		// Act
		deserialized := SubCategory{}.DeserializeEvent(serialized)

		// Assert
		require.Equal(t, event, deserialized)
	}
}
