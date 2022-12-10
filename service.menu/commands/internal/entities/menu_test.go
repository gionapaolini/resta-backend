package entities

import (
	"testing"

	"github.com/Resta-Inc/resta/pkg/events"
	"github.com/Resta-Inc/resta/pkg/eventutils2"
	"github.com/Resta-Inc/resta/pkg/resources"
	"github.com/Resta-Inc/resta/pkg/utils"
	"github.com/stretchr/testify/require"
)

func Test_CreateMenu(t *testing.T) {
	// Act
	menu := NewMenu()

	// Assert
	latestEvent := menu.Events[0]
	require.True(t, menu.IsNew())
	require.Equal(t, resources.DefaultMenuName("en"), menu.GetName())
	require.Len(t, menu.Events, 1)
	require.IsType(t, events.MenuCreated{}, latestEvent)
	require.Equal(t, utils.Time.Now(), latestEvent.GetTimeStamp())
	require.Equal(t, menu.ID, latestEvent.GetEntityID())
	require.False(t, menu.IsDeleted)
}

func Test_EnableMenu(t *testing.T) {
	// Arrange
	menu := NewMenu()

	// Act
	menu.Enable()

	// Assert
	latestEvent := menu.Events[len(menu.Events)-1]
	require.True(t, menu.IsEnabled())
	require.IsType(t, events.MenuEnabled{}, latestEvent)
}

func Test_DisableMenu(t *testing.T) {
	// Arrange
	menu := NewMenu()
	menu.Enable()

	// Act
	menu.Disable()

	// Assert
	latestEvent := menu.Events[len(menu.Events)-1]
	require.False(t, menu.IsEnabled())
	require.IsType(t, events.MenuDisabled{}, latestEvent)
}

func Test_ChangeMenuName(t *testing.T) {
	// Arrange
	menu := NewMenu()
	newName := "NewMenuName"

	// Act
	menu.ChangeName(newName)

	// Assert
	latestEvent := menu.Events[len(menu.Events)-1]
	require.Equal(t, newName, menu.GetName())
	require.IsType(t, events.MenuNameChanged{}, latestEvent)
}

func Test_AddCategory(t *testing.T) {
	// Arrange
	categoryID := utils.GenerateNewUUID()
	menu := NewMenu()

	// Act
	menu.AddCategory(categoryID)

	// Assert
	latestEvent := menu.Events[len(menu.Events)-1]
	require.Contains(t, menu.GetCategoriesIDs(), categoryID)
	require.IsType(t, events.CategoryAddedToMenu{}, latestEvent)
}

func Test_DeserializeMenuEvent(t *testing.T) {
	// Arrange
	events := []eventutils2.IEvent{
		events.MenuCreated{
			EventInfo: eventutils2.NewEventInfo(utils.GenerateNewUUID()),
		},
		events.MenuEnabled{
			EventInfo: eventutils2.NewEventInfo(utils.GenerateNewUUID()),
		},
		events.MenuDisabled{
			EventInfo: eventutils2.NewEventInfo(utils.GenerateNewUUID()),
		},
		events.MenuNameChanged{
			EventInfo: eventutils2.NewEventInfo(utils.GenerateNewUUID()),
		},
		events.CategoryAddedToMenu{
			EventInfo: eventutils2.NewEventInfo(utils.GenerateNewUUID()),
		},
	}

	for _, event := range events {
		serialized := eventutils2.SerializedEvent(event)

		// Act
		deserialized := NewMenu().DeserializeEvent(serialized)

		// Assert
		require.Equal(t, event, deserialized)
	}
}
