package entities

import (
	"testing"
	"time"

	"github.com/Resta-Inc/resta/pkg/events"
	"github.com/Resta-Inc/resta/pkg/eventutils"
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
	require.True(t, menuItem.IsNew())
	require.Equal(t, resources.DefaultMenuItemName("en"), menuItem.GetName())
	require.Len(t, menuItem.Events, 1)
	require.IsType(t, events.MenuItemCreated{}, latestEvent)
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
	require.IsType(t, events.MenuItemNameChanged{}, latestEvent)
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
	require.IsType(t, events.MenuItemEstimatedPreparationTimeChanged{}, latestEvent)
}

func Test_DeserializeMenuItemEvent(t *testing.T) {
	// Arrange
	events := []eventutils.IEvent{
		events.MenuItemCreated{
			EventInfo: eventutils.NewEventInfo(utils.GenerateNewUUID()),
		},
		events.MenuItemNameChanged{
			EventInfo: eventutils.NewEventInfo(utils.GenerateNewUUID()),
		},
		events.MenuItemEstimatedPreparationTimeChanged{
			EventInfo: eventutils.NewEventInfo(utils.GenerateNewUUID()),
		},
	}

	for _, event := range events {
		serialized := eventutils.SerializedEvent(event)

		// Act
		deserialized := MenuItem{}.DeserializeEvent(serialized)

		// Assert
		require.Equal(t, event, deserialized)
	}
}
