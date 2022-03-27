package eventutils

import (
	"testing"

	"github.com/Resta-Inc/resta/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestGetEntity(t *testing.T) {
	// Arrange
	entity := NewTestEntity()
	entity = CommitEvents(entity).(TestEntity)
	mockEventStore := new(MockEventStore)

	returnedEvents := mapToReturnedEvents(entity.GetCommittedEvents())

	mockEventStore.
		On("GetAllEventsByStreamName", GetStreamName(entity)).
		Return(returnedEvents, nil)

	repo := NewEntityRepository(mockEventStore)

	//Act
	foundEntity, err := repo.GetEntity(EmptyTestEntity(), entity.GetID())

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
		On("SaveEventsToNewStream", GetStreamName(entity), entity.GetLatestEvents()).
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
	entity = CommitEvents(entity).(TestEntity)
	entity.ChangeName("NewName")
	mockEventStore := new(MockEventStore)

	mockEventStore.
		On("SaveEventsToExistingStream", GetStreamName(entity), entity.GetLatestEvents()).
		Return(nil, nil)

	repo := NewEntityRepository(mockEventStore)

	//Act
	repo.SaveEntity(entity)

	//Assert
	mockEventStore.AssertExpectations(t)
}

func mapToReturnedEvents(events []IEvent) []ReturnedEvent {
	var returnedEvents []ReturnedEvent
	for _, event := range events {
		returnedEvents = append(returnedEvents, ReturnedEvent{
			Data: utils.SerializeObject(event),
		})
	}
	return returnedEvents
}
