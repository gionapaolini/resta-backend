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

func TestHandleMenuDisabledMessage(t *testing.T) {
	// Arrange
	menuID := utils.GenerateNewUUID()

	menuEnabledEvent := events.MenuDisabled{
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
		EntityEventInfo: eventutils.NewEntityEventInfo(menuID),
		NewName:         "NewMenuName",
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
		EntityEventInfo: eventutils.NewEntityEventInfo(categoryID),
		Name:            "TestCategoryName",
		ParentMenuID:    menuID,
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
		EntityEventInfo: eventutils.NewEntityEventInfo(menuID),
		CategoryID:      categoryID,
	}

	incomingMessage := &esdb.SubscriptionEvent{
		EventAppeared: &esdb.ResolvedEvent{
			Event: &esdb.RecordedEvent{
				Data: utils.SerializeObject(categoryAddedToMenuEvent),
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
		EntityEventInfo: eventutils.NewEntityEventInfo(categoryID),
		NewName:         "NewMenuName",
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
		EntityEventInfo:  eventutils.NewEntityEventInfo(subCategoryID),
		Name:             "TestSubCategoryName",
		ParentCategoryID: categoryID,
	}

	incomingMessage := &esdb.SubscriptionEvent{
		EventAppeared: &esdb.ResolvedEvent{
			Event: &esdb.RecordedEvent{
				Data: utils.SerializeObject(subCategoryCreatedEvent),
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
		EntityEventInfo: eventutils.NewEntityEventInfo(categoryID),
		SubCategoryID:   subCategoryID,
	}

	incomingMessage := &esdb.SubscriptionEvent{
		EventAppeared: &esdb.ResolvedEvent{
			Event: &esdb.RecordedEvent{
				Data: utils.SerializeObject(subCategoryAddedToCategoryEvent),
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
		EntityEventInfo:     eventutils.NewEntityEventInfo(menuItemID),
		Name:                "TestMenuItemName",
		ParentSubCategoryID: subCategoryID,
	}

	incomingMessage := &esdb.SubscriptionEvent{
		EventAppeared: &esdb.ResolvedEvent{
			Event: &esdb.RecordedEvent{
				Data: utils.SerializeObject(subCategoryCreatedEvent),
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
		EntityEventInfo: eventutils.NewEntityEventInfo(subCategoryID),
		MenuItemID:      menuItemID,
	}

	incomingMessage := &esdb.SubscriptionEvent{
		EventAppeared: &esdb.ResolvedEvent{
			Event: &esdb.RecordedEvent{
				Data: utils.SerializeObject(subCategoryAddedToCategoryEvent),
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
