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
	category := NewCategory()

	// Assert
	require.Equal(t, resources.DefaultCategoryName("en"), category.GetName())
	require.Equal(t, resources.DefaultCategoryImageUrl(), category.GetImageURL())
	require.Len(t, category.GetAllEvents(), 1)
	require.Len(t, category.GetCommittedEvents(), 0)
	require.Len(t, category.GetLatestEvents(), 1)
	require.IsType(t, events.CategoryCreated{}, category.GetLatestEvents()[0])
	require.Equal(t, utils.Time.Now(), category.GetLatestEvents()[0].GetDateTime())
	require.Equal(t, category.ID, category.GetLatestEvents()[0].(eventutils.IEntityEvent).GetEntityID())
	require.False(t, category.IsDeleted)
}
