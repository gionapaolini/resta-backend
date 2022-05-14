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
	require.Equal(t, resources.DefaultMenuItemName("en"), menuItem.GetName())
	require.Len(t, menuItem.GetAllEvents(), 1)
	require.Len(t, menuItem.GetCommittedEvents(), 0)
	require.Len(t, menuItem.GetLatestEvents(), 1)
	require.IsType(t, events.MenuItemCreated{}, menuItem.GetLatestEvents()[0])
	require.Equal(t, utils.Time.Now(), menuItem.GetLatestEvents()[0].GetDateTime())
	require.Equal(t, menuItem.ID, menuItem.GetLatestEvents()[0].(eventutils.IEntityEvent).GetEntityID())
	require.False(t, menuItem.IsDeleted)
}

func TestChangeMenuItemName(t *testing.T) {
	// Arrange
	menuItemID := utils.GenerateNewUUID()
	menuItem := NewMenuItem(menuItemID)

	// Act
	menuItem = menuItem.ChangeName("NewName")

	// Assert
	require.Equal(t, "NewName", menuItem.GetName())
	require.IsType(t, events.MenuItemNameChanged{}, menuItem.GetLatestEvents()[len(menuItem.GetLatestEvents())-1])
}

func TestChangeMenuItemEstimatedPreparationTime(t *testing.T) {
	// Arrange
	menuItemID := utils.GenerateNewUUID()
	menuItem := NewMenuItem(menuItemID)

	// Act
	menuItem = menuItem.ChangeEstimatedPreparationTime(10 * time.Minute)

	// Assert
	require.Equal(t, 10*time.Minute, menuItem.GetEstimatedPreparationtime())
	require.IsType(t, events.MenuItemEstimatedPreparationTimeChanged{}, menuItem.GetLatestEvents()[len(menuItem.GetLatestEvents())-1])
}
