package entities

import (
	"testing"

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
	require.Equal(t, resources.DefaultMenuItemName("en"), menuItem.GetName())
	require.Len(t, menuItem.GetAllEvents(), 1)
	require.Len(t, menuItem.GetCommittedEvents(), 0)
	require.Len(t, menuItem.GetLatestEvents(), 1)
	require.IsType(t, events.MenuItemCreated{}, menuItem.GetLatestEvents()[0])
	require.Equal(t, utils.Time.Now(), menuItem.GetLatestEvents()[0].GetDateTime())
	require.Equal(t, menuItem.ID, menuItem.GetLatestEvents()[0].(eventutils.IEntityEvent).GetEntityID())
	require.False(t, menuItem.IsDeleted)
}
