package internal

import (
	"testing"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/Resta-Inc/resta/menu/commands/internal/entities"
	"github.com/Resta-Inc/resta/pkg/events"
	"github.com/Resta-Inc/resta/pkg/eventutils"
	"github.com/Resta-Inc/resta/pkg/utils"
	"github.com/stretchr/testify/mock"
	"golang.org/x/exp/slices"
)

func TestHandleCategoryCreatedMessage(t *testing.T) {
	// Arrange
	categoryID := utils.GenerateNewUUID()
	menu := entities.NewMenu()

	categoryCreatedEvent := events.CategoryCreated{
		EntityEventInfo: eventutils.NewEntityEventInfo(categoryID),
		Name:            "TestCategoryName",
		ParentMenuID:    menu.ID,
	}

	incomingMessage := &esdb.SubscriptionEvent{
		EventAppeared: &esdb.ResolvedEvent{
			Event: &esdb.RecordedEvent{
				Data: utils.SerializeObject(categoryCreatedEvent),
			},
		},
		SubscriptionDropped: &esdb.SubscriptionDropped{},
		CheckPointReached:   &esdb.Position{},
	}

	mockEntityRepository := new(eventutils.MockEntityRepository)
	mockEntityRepository.
		On("GetEntity", entities.EmptyMenu(), menu.ID).
		Return(menu, nil)

	mockEntityRepository.
		On("SaveEntity", mock.MatchedBy(
			func(menu entities.Menu) bool {
				return slices.Contains(menu.GetCategoriesIDs(), categoryID)
			},
		)).
		Return(nil)

	eventHandler := NewEventHandler(mockEntityRepository)

	// Act
	eventHandler.HandleCategoryCreated(incomingMessage)

	// Assert
	mockEntityRepository.AssertExpectations(t)
}
