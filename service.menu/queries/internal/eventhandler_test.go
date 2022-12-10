package internal

import (
	"testing"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/Resta-Inc/resta/pkg/events"
	"github.com/Resta-Inc/resta/pkg/eventutils2"
	"github.com/Resta-Inc/resta/pkg/utils"
)

func TestHandleMenuCreatedMessage(t *testing.T) {
	// Arrange
	menuID := utils.GenerateNewUUID()

	menuCreatedEvent := events.MenuCreated{
		EventInfo: eventutils2.NewEventInfo(menuID),
		Name:      "TestMenuName",
	}

	serializedEvent := eventutils2.SerializedEvent(menuCreatedEvent)

	incomingMessage := &esdb.SubscriptionEvent{
		EventAppeared: &esdb.ResolvedEvent{
			Event: &esdb.RecordedEvent{
				EventID:   serializedEvent.ID,
				EventType: serializedEvent.Name,
				Data:      serializedEvent.Data,
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
		EventInfo: eventutils2.NewEventInfo(menuID),
	}

	serializedEvent := eventutils2.SerializedEvent(menuEnabledEvent)

	incomingMessage := &esdb.SubscriptionEvent{
		EventAppeared: &esdb.ResolvedEvent{
			Event: &esdb.RecordedEvent{
				EventID:   serializedEvent.ID,
				EventType: serializedEvent.Name,
				Data:      serializedEvent.Data,
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

func TestHandleMenuDisabledMessage(t *testing.T) {
	// Arrange
	menuID := utils.GenerateNewUUID()

	menuEnabledEvent := events.MenuDisabled{
		EventInfo: eventutils2.NewEventInfo(menuID),
	}

	serializedEvent := eventutils2.SerializedEvent(menuEnabledEvent)

	incomingMessage := &esdb.SubscriptionEvent{
		EventAppeared: &esdb.ResolvedEvent{
			Event: &esdb.RecordedEvent{
				EventID:   serializedEvent.ID,
				EventType: serializedEvent.Name,
				Data:      serializedEvent.Data,
			},
		},
		SubscriptionDropped: &esdb.SubscriptionDropped{},
		CheckPointReached:   &esdb.Position{},
	}

	mockMenuRepository := new(MockMenuRepository)
	mockMenuRepository.
		On("DisableMenu", menuID).
		Return(nil)

	eventHandler := NewMenuEventHandler(mockMenuRepository)

	// Act
	eventHandler.HandleMenuDisabled(incomingMessage)

	// Assert
	mockMenuRepository.AssertExpectations(t)
}

func TestHandleMenuNameChangedMessage(t *testing.T) {
	// Arrange
	menuID := utils.GenerateNewUUID()

	menuEnabledEvent := events.MenuNameChanged{
		EventInfo: eventutils2.NewEventInfo(menuID),
		NewName:   "NewMenuName",
	}

	serializedEvent := eventutils2.SerializedEvent(menuEnabledEvent)

	incomingMessage := &esdb.SubscriptionEvent{
		EventAppeared: &esdb.ResolvedEvent{
			Event: &esdb.RecordedEvent{
				EventID:   serializedEvent.ID,
				EventType: serializedEvent.Name,
				Data:      serializedEvent.Data,
			},
		},
		SubscriptionDropped: &esdb.SubscriptionDropped{},
		CheckPointReached:   &esdb.Position{},
	}

	mockMenuRepository := new(MockMenuRepository)
	mockMenuRepository.
		On("ChangeMenuName", menuID, menuEnabledEvent.NewName).
		Return(nil)

	eventHandler := NewMenuEventHandler(mockMenuRepository)

	// Act
	eventHandler.HandleMenuNameChanged(incomingMessage)

	// Assert
	mockMenuRepository.AssertExpectations(t)
}

func TestHandleCategoryCreatedMessage(t *testing.T) {
	// Arrange
	menuID, categoryID := utils.GenerateNewUUID(), utils.GenerateNewUUID()

	categoryCreatedEvent := events.CategoryCreated{
		EventInfo:    eventutils2.NewEventInfo(categoryID),
		Name:         "TestCategoryName",
		ParentMenuID: menuID,
	}

	serializedEvent := eventutils2.SerializedEvent(categoryCreatedEvent)

	incomingMessage := &esdb.SubscriptionEvent{
		EventAppeared: &esdb.ResolvedEvent{
			Event: &esdb.RecordedEvent{
				EventID:   serializedEvent.ID,
				EventType: serializedEvent.Name,
				Data:      serializedEvent.Data,
			},
		},
		SubscriptionDropped: &esdb.SubscriptionDropped{},
		CheckPointReached:   &esdb.Position{},
	}

	mockMenuRepository := new(MockMenuRepository)
	mockMenuRepository.
		On("CreateCategory",
			categoryID,
			categoryCreatedEvent.Name).
		Return(nil)

	eventHandler := NewMenuEventHandler(mockMenuRepository)

	// Act
	eventHandler.HandleCategoryCreated(incomingMessage)

	// Assert
	mockMenuRepository.AssertExpectations(t)
}

func TestHandleCategoryAddedToMenuMessage(t *testing.T) {
	// Arrange
	menuID, categoryID := utils.GenerateNewUUID(), utils.GenerateNewUUID()

	categoryAddedToMenuEvent := events.CategoryAddedToMenu{
		EventInfo:  eventutils2.NewEventInfo(menuID),
		CategoryID: categoryID,
	}

	serializedEvent := eventutils2.SerializedEvent(categoryAddedToMenuEvent)

	incomingMessage := &esdb.SubscriptionEvent{
		EventAppeared: &esdb.ResolvedEvent{
			Event: &esdb.RecordedEvent{
				EventID:   serializedEvent.ID,
				EventType: serializedEvent.Name,
				Data:      serializedEvent.Data,
			},
		},
		SubscriptionDropped: &esdb.SubscriptionDropped{},
		CheckPointReached:   &esdb.Position{},
	}

	mockMenuRepository := new(MockMenuRepository)
	mockMenuRepository.
		On("AddCategoryToMenu", menuID, categoryID).
		Return(nil)

	eventHandler := NewMenuEventHandler(mockMenuRepository)

	// Act
	eventHandler.HandleCategoryAddedToMenu(incomingMessage)

	// Assert
	mockMenuRepository.AssertExpectations(t)
}

func TestHandleCategoryNameChangedMessage(t *testing.T) {
	// Arrange
	categoryID := utils.GenerateNewUUID()

	menuEnabledEvent := events.CategoryNameChanged{
		EventInfo: eventutils2.NewEventInfo(categoryID),
		NewName:   "NewMenuName",
	}

	serializedEvent := eventutils2.SerializedEvent(menuEnabledEvent)

	incomingMessage := &esdb.SubscriptionEvent{
		EventAppeared: &esdb.ResolvedEvent{
			Event: &esdb.RecordedEvent{
				EventID:   serializedEvent.ID,
				EventType: serializedEvent.Name,
				Data:      serializedEvent.Data,
			},
		},
		SubscriptionDropped: &esdb.SubscriptionDropped{},
		CheckPointReached:   &esdb.Position{},
	}

	mockMenuRepository := new(MockMenuRepository)
	mockMenuRepository.
		On("ChangeCategoryName", categoryID, menuEnabledEvent.NewName).
		Return(nil)

	eventHandler := NewMenuEventHandler(mockMenuRepository)

	// Act
	eventHandler.HandleCategoryNameChanged(incomingMessage)

	// Assert
	mockMenuRepository.AssertExpectations(t)
}
func TestHandleSubCategoryCreatedMessage(t *testing.T) {
	// Arrange
	categoryID, subCategoryID := utils.GenerateNewUUID(), utils.GenerateNewUUID()

	subCategoryCreatedEvent := events.SubCategoryCreated{
		EventInfo:        eventutils2.NewEventInfo(subCategoryID),
		Name:             "TestSubCategoryName",
		ParentCategoryID: categoryID,
	}

	serializedEvent := eventutils2.SerializedEvent(subCategoryCreatedEvent)

	incomingMessage := &esdb.SubscriptionEvent{
		EventAppeared: &esdb.ResolvedEvent{
			Event: &esdb.RecordedEvent{
				EventID:   serializedEvent.ID,
				EventType: serializedEvent.Name,
				Data:      serializedEvent.Data,
			},
		},
		SubscriptionDropped: &esdb.SubscriptionDropped{},
		CheckPointReached:   &esdb.Position{},
	}

	mockMenuRepository := new(MockMenuRepository)
	mockMenuRepository.
		On("CreateSubCategory",
			subCategoryID,
			subCategoryCreatedEvent.Name).
		Return(nil)

	eventHandler := NewMenuEventHandler(mockMenuRepository)

	// Act
	eventHandler.HandleSubCategoryCreated(incomingMessage)

	// Assert
	mockMenuRepository.AssertExpectations(t)
}

func TestHandleSubCategoryAddetToCategory(t *testing.T) {
	// Arrange
	categoryID, subCategoryID := utils.GenerateNewUUID(), utils.GenerateNewUUID()

	subCategoryAddedToCategoryEvent := events.SubCategoryAddedToCategory{
		EventInfo:     eventutils2.NewEventInfo(categoryID),
		SubCategoryID: subCategoryID,
	}

	serializedEvent := eventutils2.SerializedEvent(subCategoryAddedToCategoryEvent)

	incomingMessage := &esdb.SubscriptionEvent{
		EventAppeared: &esdb.ResolvedEvent{
			Event: &esdb.RecordedEvent{
				EventID:   serializedEvent.ID,
				EventType: serializedEvent.Name,
				Data:      serializedEvent.Data,
			},
		},
		SubscriptionDropped: &esdb.SubscriptionDropped{},
		CheckPointReached:   &esdb.Position{},
	}

	mockMenuRepository := new(MockMenuRepository)
	mockMenuRepository.
		On("AddSubCategoryToCategory",
			categoryID,
			subCategoryID).
		Return(nil)

	eventHandler := NewMenuEventHandler(mockMenuRepository)

	// Act
	eventHandler.HandleSubCategoryAddedToCategory(incomingMessage)

	// Assert
	mockMenuRepository.AssertExpectations(t)
}

func TestHandleMenuItemCreatedMessage(t *testing.T) {
	// Arrange
	subCategoryID, menuItemID := utils.GenerateNewUUID(), utils.GenerateNewUUID()

	subCategoryCreatedEvent := events.MenuItemCreated{
		EventInfo:           eventutils2.NewEventInfo(menuItemID),
		Name:                "TestMenuItemName",
		ParentSubCategoryID: subCategoryID,
	}

	serializedEvent := eventutils2.SerializedEvent(subCategoryCreatedEvent)

	incomingMessage := &esdb.SubscriptionEvent{
		EventAppeared: &esdb.ResolvedEvent{
			Event: &esdb.RecordedEvent{
				EventID:   serializedEvent.ID,
				EventType: serializedEvent.Name,
				Data:      serializedEvent.Data,
			},
		},
		SubscriptionDropped: &esdb.SubscriptionDropped{},
		CheckPointReached:   &esdb.Position{},
	}

	mockMenuRepository := new(MockMenuRepository)
	mockMenuRepository.
		On("CreateMenuItem",
			menuItemID,
			subCategoryCreatedEvent.Name).
		Return(nil)

	eventHandler := NewMenuEventHandler(mockMenuRepository)

	// Act
	eventHandler.HandleMenuItemCreated(incomingMessage)

	// Assert
	mockMenuRepository.AssertExpectations(t)
}

func TestHandleMenuItemAddetToCategory(t *testing.T) {
	// Arrange
	subCategoryID, menuItemID := utils.GenerateNewUUID(), utils.GenerateNewUUID()

	subCategoryAddedToCategoryEvent := events.MenuItemAddedToSubCategory{
		EventInfo:  eventutils2.NewEventInfo(subCategoryID),
		MenuItemID: menuItemID,
	}

	serializedEvent := eventutils2.SerializedEvent(subCategoryAddedToCategoryEvent)

	incomingMessage := &esdb.SubscriptionEvent{
		EventAppeared: &esdb.ResolvedEvent{
			Event: &esdb.RecordedEvent{
				EventID:   serializedEvent.ID,
				EventType: serializedEvent.Name,
				Data:      serializedEvent.Data,
			},
		},
		SubscriptionDropped: &esdb.SubscriptionDropped{},
		CheckPointReached:   &esdb.Position{},
	}

	mockMenuRepository := new(MockMenuRepository)
	mockMenuRepository.
		On("AddMenuItemToSubCategory",
			subCategoryID,
			menuItemID).
		Return(nil)

	eventHandler := NewMenuEventHandler(mockMenuRepository)

	// Act
	eventHandler.HandleMenuItemAddedToSubCategory(incomingMessage)

	// Assert
	mockMenuRepository.AssertExpectations(t)
}
