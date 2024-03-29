package entities

import (
	"testing"

	"github.com/Resta-Inc/resta/pkg/events"
	"github.com/Resta-Inc/resta/pkg/eventutils"
	"github.com/Resta-Inc/resta/pkg/resources"
	"github.com/Resta-Inc/resta/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestCreateCategory(t *testing.T) {
	// Act
	menuID := utils.GenerateNewUUID()
	category := NewCategory(menuID)

	// Assert
	latestEvent := category.Events[len(category.Events)-1]
	require.True(t, category.IsNew())
	require.Equal(t, resources.DefaultCategoryName("en"), category.GetName())
	require.Len(t, category.Events, 1)
	require.IsType(t, events.CategoryCreated{}, latestEvent)
	require.Equal(t, utils.Time.Now(), latestEvent.GetTimeStamp())
	require.Equal(t, category.ID, latestEvent.GetEntityID())
	require.False(t, category.IsDeleted)
}

func TestChangeCategoryName(t *testing.T) {
	// Arrange
	menuID := utils.GenerateNewUUID()
	category := NewCategory(menuID)
	newName := "New name"

	// Act
	category.ChangeName(newName)

	// Assert
	latestEvent := category.Events[len(category.Events)-1]
	require.Equal(t, newName, category.GetName())
	require.IsType(t, events.CategoryNameChanged{}, latestEvent)
}

func Test_AddSubCategory(t *testing.T) {
	// Arrange
	subCategoryID := utils.GenerateNewUUID()
	category := NewCategory(utils.GenerateNewUUID())

	// Act
	category.AddSubCategory(subCategoryID)

	// Assert
	latestEvent := category.Events[len(category.Events)-1]
	require.Contains(t, category.GetSubCategoriesIDs(), subCategoryID)
	require.IsType(t, events.SubCategoryAddedToCategory{}, latestEvent)
}

func Test_DeserializeCategoryEvent(t *testing.T) {
	// Arrange
	events := []eventutils.IEvent{
		events.CategoryCreated{
			EventInfo: eventutils.NewEventInfo(utils.GenerateNewUUID()),
		},
		events.CategoryNameChanged{
			EventInfo: eventutils.NewEventInfo(utils.GenerateNewUUID()),
		},
		events.SubCategoryAddedToCategory{
			EventInfo: eventutils.NewEventInfo(utils.GenerateNewUUID()),
		},
	}

	for _, event := range events {
		serialized := eventutils.SerializedEvent(event)

		// Act
		deserialized := Category{}.DeserializeEvent(serialized)

		// Assert
		require.Equal(t, event, deserialized)
	}
}
