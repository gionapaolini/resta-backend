package entities

import (
	"testing"

	"github.com/Resta-Inc/resta/pkg/events"
	"github.com/Resta-Inc/resta/pkg/eventutils"
	"github.com/Resta-Inc/resta/pkg/resources"
	"github.com/Resta-Inc/resta/pkg/utils"
	"github.com/stretchr/testify/require"
)

func Test_CreateMenu(t *testing.T) {
	// Act
	menu := NewMenu()

	// Assert
	require.Equal(t, resources.DefaultMenuName("en"), menu.GetName())
	require.Len(t, menu.GetAllEvents(), 1)
	require.Len(t, menu.GetCommittedEvents(), 0)
	require.Len(t, menu.GetLatestEvents(), 1)
	require.IsType(t, events.MenuCreated{}, menu.GetLatestEvents()[0])
	require.Equal(t, utils.Time.Now(), menu.GetLatestEvents()[0].GetDateTime())
	require.Equal(t, menu.ID, menu.GetLatestEvents()[0].(eventutils.IEntityEvent).GetEntityID())
	require.False(t, menu.IsDeleted)
}

func Test_EnableMenu(t *testing.T) {
	// Arrange
	menu := NewMenu()

	// Act
	menu = menu.Enable()

	// Assert
	require.True(t, menu.IsEnabled())
	require.IsType(t, events.MenuEnabled{}, menu.GetLatestEvents()[1])
}

func Test_DisableMenu(t *testing.T) {
	// Arrange
	menu := NewMenu()
	menu = menu.Enable()

	// Act
	menu = menu.Disable()

	// Assert
	require.False(t, menu.IsEnabled())
	require.IsType(t, events.MenuDisabled{}, menu.GetLatestEvents()[2])
}

func Test_ChangeMenuName(t *testing.T) {
	// Arrange
	menu := NewMenu()
	newName := "NewMenuName"

	// Act
	menu = menu.ChangeName(newName)

	// Assert
	require.Equal(t, newName, menu.GetName())
	require.IsType(t, events.MenuNameChanged{}, menu.GetLatestEvents()[1])
}

func Test_AddCategory(t *testing.T) {
	// Arrange
	categoryID := utils.GenerateNewUUID()
	menu := NewMenu()

	// Act
	menu = menu.AddCategory(categoryID)

	// Assert
	require.Contains(t, menu.GetCategoriesIDs(), categoryID)
	require.IsType(t, events.CategoryAddedToMenu{}, menu.GetLatestEvents()[1])
}

func Test_DeserializeMenuEvent(t *testing.T) {
	// Arrange
	events := []eventutils.IEvent{
		events.MenuCreated{
			EntityEventInfo: eventutils.NewEntityEventInfo(utils.GenerateNewUUID()),
		},
		events.MenuEnabled{
			EntityEventInfo: eventutils.NewEntityEventInfo(utils.GenerateNewUUID()),
		},
		events.MenuDisabled{
			EntityEventInfo: eventutils.NewEntityEventInfo(utils.GenerateNewUUID()),
		},
		events.MenuNameChanged{
			EntityEventInfo: eventutils.NewEntityEventInfo(utils.GenerateNewUUID()),
		},
		events.CategoryAddedToMenu{
			EntityEventInfo: eventutils.NewEntityEventInfo(utils.GenerateNewUUID()),
		},
	}

	for _, event := range events {
		serialized := utils.SerializeObject(event)

		// Act
		deserialized := EmptyMenu().DeserializeEvent(serialized)

		// Assert
		require.Equal(t, event, deserialized)
	}
}
