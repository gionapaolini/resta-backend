package internal

import (
	"testing"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/Resta-Inc/resta/pkg/events"
	"github.com/Resta-Inc/resta/pkg/eventutils"
	"github.com/Resta-Inc/resta/pkg/utils"
)

func TestHandleMenuCreatedMessage(t *testing.T) {
	// Arrange
	menuID := utils.GenerateNewUUID()

	menuCreatedEvent := events.MenuCreated{
		EntityEventInfo: eventutils.NewEntityEventInfo(menuID),
		Name:            "TestMenuName",
	}

	incomingMessage := &esdb.SubscriptionEvent{
		EventAppeared: &esdb.ResolvedEvent{
			Event: &esdb.RecordedEvent{
				Data: utils.SerializeObject(menuCreatedEvent),
			},
		},
		SubscriptionDropped: &esdb.SubscriptionDropped{},
		CheckPointReached:   &esdb.Position{},
	}

	mockMenuRepository := new(MockMenuRepository)
	mockMenuRepository.
		On("CreateMenu", menuID, menuCreatedEvent.Name).
		Return(nil)

	eventHandler := NewMenuEventHandler(mockMenuRepository)

	// Act
	eventHandler.HandleMenuCreated(incomingMessage)

	// Assert
	mockMenuRepository.AssertExpectations(t)
}

func TestHandleMenuEnabledMessage(t *testing.T) {
	// Arrange
	menuID := utils.GenerateNewUUID()

	menuEnabledEvent := events.MenuEnabled{
		EntityEventInfo: eventutils.NewEntityEventInfo(menuID),
	}

	incomingMessage := &esdb.SubscriptionEvent{
		EventAppeared: &esdb.ResolvedEvent{
			Event: &esdb.RecordedEvent{
				Data: utils.SerializeObject(menuEnabledEvent),
			},
		},
		SubscriptionDropped: &esdb.SubscriptionDropped{},
		CheckPointReached:   &esdb.Position{},
	}

	mockMenuRepository := new(MockMenuRepository)
	mockMenuRepository.
		On("EnableMenu", menuID).
		Return(nil)

	eventHandler := NewMenuEventHandler(mockMenuRepository)

	// Act
	eventHandler.HandleMenuEnabled(incomingMessage)

	// Assert
	mockMenuRepository.AssertExpectations(t)
}
