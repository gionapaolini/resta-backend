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
	require.Equal(t, resources.DefaultCategoryName("en"), category.GetName())
	require.Len(t, category.GetAllEvents(), 1)
	require.Len(t, category.GetCommittedEvents(), 0)
	require.Len(t, category.GetLatestEvents(), 1)
	require.IsType(t, events.CategoryCreated{}, category.GetLatestEvents()[0])
	require.Equal(t, utils.Time.Now(), category.GetLatestEvents()[0].GetDateTime())
	require.Equal(t, category.ID, category.GetLatestEvents()[0].(eventutils.IEntityEvent).GetEntityID())
	require.False(t, category.IsDeleted)
}

func TestChangeCategoryName(t *testing.T) {
	// Arrange
	menuID := utils.GenerateNewUUID()
	category := NewCategory(menuID)
	newName := "New name"

	// Act
	category = category.ChangeName(newName)

	// Assert
	require.Equal(t, newName, category.GetName())
	require.IsType(t, events.CategoryNameChanged{}, category.GetLatestEvents()[1])
}

func Test_AddSubCategory(t *testing.T) {
	// Arrange
	subCategoryID := utils.GenerateNewUUID()
	category := NewCategory(utils.GenerateNewUUID())

	// Act
	category = category.AddSubCategory(subCategoryID)

	// Assert
	require.Contains(t, category.GetSubCategoriesIDs(), subCategoryID)
	require.IsType(t, events.SubCategoryAddedToCategory{}, category.GetLatestEvents()[1])
}

func Test_DeserializeCategoryEvent(t *testing.T) {
	// Arrange
	events := []eventutils.IEvent{
		events.CategoryCreated{
			EntityEventInfo: eventutils.NewEntityEventInfo(utils.GenerateNewUUID()),
		},
		events.CategoryNameChanged{
			EntityEventInfo: eventutils.NewEntityEventInfo(utils.GenerateNewUUID()),
		},
	}

	for _, event := range events {
		serialized := utils.SerializeObject(event)

		// Act
		deserialized := EmptyCategory().DeserializeEvent(serialized)

		// Assert
		require.Equal(t, event, deserialized)
	}
}
