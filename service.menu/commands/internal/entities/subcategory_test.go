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
	require.Equal(t, resources.DefaultSubCategoryName("en"), subCategory.GetName())
	require.Len(t, subCategory.GetAllEvents(), 1)
	require.Len(t, subCategory.GetCommittedEvents(), 0)
	require.Len(t, subCategory.GetLatestEvents(), 1)
	require.IsType(t, events.SubCategoryCreated{}, subCategory.GetLatestEvents()[0])
	require.Equal(t, utils.Time.Now(), subCategory.GetLatestEvents()[0].GetDateTime())
	require.Equal(t, subCategory.ID, subCategory.GetLatestEvents()[0].(eventutils.IEntityEvent).GetEntityID())
	require.False(t, subCategory.IsDeleted)
}

func TestChangeSubCategoryName(t *testing.T) {
	// Arrange
	categoryID := utils.GenerateNewUUID()
	subCategory := NewSubCategory(categoryID)
	newName := "New name"
	// Act
	subCategory = subCategory.ChangeName(newName)

	// Assert
	require.Equal(t, newName, subCategory.GetName())
	require.IsType(t, events.SubCategoryNameChanged{}, subCategory.GetLatestEvents()[1])
}

func Test_DeserializeSubCategoryEvent(t *testing.T) {
	// Arrange
	events := []eventutils.IEvent{
		events.SubCategoryCreated{
			EntityEventInfo: eventutils.NewEntityEventInfo(utils.GenerateNewUUID()),
		},
		events.SubCategoryNameChanged{
			EntityEventInfo: eventutils.NewEntityEventInfo(utils.GenerateNewUUID()),
		},
	}

	for _, event := range events {
		serialized := utils.SerializeObject(event)

		// Act
		deserialized := EmptySubCategory().DeserializeEvent(serialized)

		// Assert
		require.Equal(t, event, deserialized)
	}
}
