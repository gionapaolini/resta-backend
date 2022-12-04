package entities2

import (
	"testing"
	"time"

	"github.com/Resta-Inc/resta/pkg/events2"
	"github.com/Resta-Inc/resta/pkg/eventutils2"
	"github.com/Resta-Inc/resta/pkg/resources"
	"github.com/Resta-Inc/resta/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestCreateMenuItem(t *testing.T) {
	// Act
	menuItemID := utils.GenerateNewUUID()
	menuItem := NewMenuItem(menuItemID)

	// Assert
	latestEvent := menuItem.Events[len(menuItem.Events)-1]
	require.Equal(t, resources.DefaultMenuItemName("en"), menuItem.GetName())
	require.Len(t, menuItem.Events, 1)
	require.IsType(t, events2.MenuItemCreated{}, latestEvent)
	require.Equal(t, utils.Time.Now(), latestEvent.GetTimeStamp())
	require.Equal(t, menuItem.ID, latestEvent.GetEntityID())
	require.False(t, menuItem.IsDeleted)
}

func TestChangeMenuItemName(t *testing.T) {
	// Arrange
	menuItemID := utils.GenerateNewUUID()
	menuItem := NewMenuItem(menuItemID)

	// Act
	menuItem.ChangeName("NewName")

	// Assert
	latestEvent := menuItem.Events[len(menuItem.Events)-1]
	require.Equal(t, "NewName", menuItem.GetName())
	require.IsType(t, events2.MenuItemNameChanged{}, latestEvent)
}

func TestChangeMenuItemEstimatedPreparationTime(t *testing.T) {
	// Arrange
	menuItemID := utils.GenerateNewUUID()
	menuItem := NewMenuItem(menuItemID)

	// Act
	menuItem.ChangeEstimatedPreparationTime(10 * time.Minute)

	// Assert
	latestEvent := menuItem.Events[len(menuItem.Events)-1]
	require.Equal(t, 10*time.Minute, menuItem.GetEstimatedPreparationtime())
	require.IsType(t, events2.MenuItemEstimatedPreparationTimeChanged{}, latestEvent)
}

func Test_DeserializeMenuItemEvent(t *testing.T) {
	// Arrange
	events := []eventutils2.IEvent{
		events2.MenuItemCreated{
			EventInfo: eventutils2.NewEventInfo(utils.GenerateNewUUID()),
		},
		events2.MenuItemNameChanged{
			EventInfo: eventutils2.NewEventInfo(utils.GenerateNewUUID()),
		},
		events2.MenuItemEstimatedPreparationTimeChanged{
			EventInfo: eventutils2.NewEventInfo(utils.GenerateNewUUID()),
		},
	}

	for _, event := range events {
		serialized := eventutils2.SerializedEvent(event)

		// Act
		deserialized := MenuItem{}.DeserializeEvent(serialized)

		// Assert
		require.Equal(t, event, deserialized)
	}
}
