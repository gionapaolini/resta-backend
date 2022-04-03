package entities

import (
	"testing"

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
	require.IsType(t, MenuCreated{}, menu.GetLatestEvents()[0])
	require.Equal(t, utils.Time.Now(), menu.GetLatestEvents()[0].GetDateTime())
	require.Equal(t, menu.ID, menu.GetLatestEvents()[0].(eventutils.IEntityEvent).GetEntityID())
	require.False(t, menu.IsDeleted)
}

func Test_DeserializeMenuEvent(t *testing.T) {
	// Arrange
	events := []eventutils.IEvent{
		MenuCreated{
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