package eventutils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetEntity(t *testing.T) {
	// Arrange
	entity := NewTestEntity()
	mockEventStore := new(MockEventStore)

	serializedEvents := serializeEvents(entity.GetEvents())

	// resetting the entity like it was saved, it will allow to compare with the retrieved one
	entity.New = false
	entity.Events = nil

	mockEventStore.
		On("GetAllEventsByStreamName", getStreamName(entity)).
		Return(serializedEvents, nil)

	repo := NewEntityRepository(mockEventStore)

	//Act
	foundEntity, err := repo.GetEntity(&TestEntity{}, entity.GetID())

	//Assert
	mockEventStore.AssertExpectations(t)
	require.NoError(t, err)
	require.Equal(t, entity, foundEntity)
}

func TestSaveNewEntity(t *testing.T) {
	// Arrange
	entity := NewTestEntity()
	mockEventStore := new(MockEventStore)

	mockEventStore.
		On("SaveEventsToNewStream", getStreamName(entity), serializeEvents(entity.Events)).
		Return(nil, nil)

	repo := NewEntityRepository(mockEventStore)

	//Act
	repo.SaveEntity(entity)

	//Assert
	mockEventStore.AssertExpectations(t)
}

func TestSaveExistingEntity(t *testing.T) {
	// Arrange
	entity := NewTestEntity()
	entity.Events = []IEvent{}
	entity.New = false
	entity.ChangeName("NewName")
	mockEventStore := new(MockEventStore)

	mockEventStore.
		On("SaveEventsToExistingStream", getStreamName(entity), serializeEvents(entity.Events)).
		Return(nil, nil)

	repo := NewEntityRepository(mockEventStore)

	//Act
	repo.SaveEntity(entity)

	//Assert
	mockEventStore.AssertExpectations(t)
}
