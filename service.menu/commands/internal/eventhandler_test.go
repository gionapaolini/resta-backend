package internal

import (
	"testing"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/Resta-Inc/resta/menu/commands/internal/entities2"
	"github.com/Resta-Inc/resta/pkg/events2"
	"github.com/Resta-Inc/resta/pkg/eventutils2"
	"github.com/Resta-Inc/resta/pkg/utils"
	"github.com/stretchr/testify/mock"
	"golang.org/x/exp/slices"
)

func TestHandleCategoryCreatedMessage(t *testing.T) {
	// Arrange
	categoryID := utils.GenerateNewUUID()
	menu := entities2.NewMenu()

	categoryCreatedEvent := events2.CategoryCreated{
		EventInfo:    eventutils2.NewEventInfo(categoryID),
		Name:         "TestCategoryName",
		ParentMenuID: menu.ID,
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

	mockEntityRepository := new(eventutils2.MockEntityRepository)
	mockEntityRepository.
		On("GetEntity", &entities2.Menu{}, menu.ID).
		Return(menu, nil)

	mockEntityRepository.
		On("SaveEntity", mock.MatchedBy(
			func(menu *entities2.Menu) bool {
				return slices.Contains(menu.GetCategoriesIDs(), categoryID)
			},
		)).
		Return(nil)

	eventHandler := NewMenuEventHandler(mockEntityRepository)

	// Act
	eventHandler.HandleCategoryCreated(incomingMessage)

	// Assert
	mockEntityRepository.AssertExpectations(t)
}

func TestHandleSubCategoryCreatedMessage(t *testing.T) {
	// Arrange
	subCategoryID := utils.GenerateNewUUID()
	category := entities2.NewCategory(utils.GenerateNewUUID())

	subCategoryCreatedEvent := events2.SubCategoryCreated{
		EventInfo:        eventutils2.NewEventInfo(subCategoryID),
		Name:             "TestCategoryName",
		ParentCategoryID: category.ID,
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

	mockEntityRepository := new(eventutils2.MockEntityRepository)
	mockEntityRepository.
		On("GetEntity", &entities2.Category{}, category.ID).
		Return(category, nil)

	mockEntityRepository.
		On("SaveEntity", mock.MatchedBy(
			func(category *entities2.Category) bool {
				return slices.Contains(category.GetSubCategoriesIDs(), subCategoryID)
			},
		)).
		Return(nil)

	eventHandler := NewMenuEventHandler(mockEntityRepository)

	// Act
	eventHandler.HandleSubCategoryCreated(incomingMessage)

	// Assert
	mockEntityRepository.AssertExpectations(t)
}

func TestHandleMenuItemCreatedMessage(t *testing.T) {
	// Arrange
	menuItemID := utils.GenerateNewUUID()
	subCategory := entities2.NewSubCategory(utils.GenerateNewUUID())

	menuItemCreatedEvent := events2.MenuItemCreated{
		EventInfo:           eventutils2.NewEventInfo(menuItemID),
		Name:                "TestMenuItemName",
		ParentSubCategoryID: subCategory.ID,
	}

	serializedEvent := eventutils2.SerializedEvent(menuItemCreatedEvent)

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

	mockEntityRepository := new(eventutils2.MockEntityRepository)
	mockEntityRepository.
		On("GetEntity", &entities2.SubCategory{}, subCategory.ID).
		Return(subCategory, nil)

	mockEntityRepository.
		On("SaveEntity", mock.MatchedBy(
			func(subCategory *entities2.SubCategory) bool {
				return slices.Contains(subCategory.GetMenuItemsIDs(), menuItemID)
			},
		)).
		Return(nil)

	eventHandler := NewMenuEventHandler(mockEntityRepository)

	// Act
	eventHandler.HandleMenuItemCreated(incomingMessage)

	// Assert
	mockEntityRepository.AssertExpectations(t)
}
