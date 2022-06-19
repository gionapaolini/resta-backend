package entities2

import (
	"testing"

	"github.com/Resta-Inc/resta/pkg/events2"
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
	require.Equal(t, resources.DefaultMenuName("en"), menu.GetName())
	require.Len(t, menu.Events, 1)
	require.IsType(t, events2.MenuCreated{}, latestEvent)
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
	require.IsType(t, events2.MenuEnabled{}, latestEvent)
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
	require.IsType(t, events2.MenuDisabled{}, latestEvent)
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
	require.IsType(t, events2.MenuNameChanged{}, latestEvent)
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
	require.IsType(t, events2.CategoryAddedToMenu{}, latestEvent)
}

func Test_DeserializeMenuEvent(t *testing.T) {
	// Arrange
	events := []eventutils2.IEvent{
		events2.MenuCreated{
			EventInfo: eventutils2.NewEventInfo(utils.GenerateNewUUID()),
		},
		events2.MenuEnabled{
			EventInfo: eventutils2.NewEventInfo(utils.GenerateNewUUID()),
		},
		events2.MenuDisabled{
			EventInfo: eventutils2.NewEventInfo(utils.GenerateNewUUID()),
		},
		events2.MenuNameChanged{
			EventInfo: eventutils2.NewEventInfo(utils.GenerateNewUUID()),
		},
		events2.CategoryAddedToMenu{
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
