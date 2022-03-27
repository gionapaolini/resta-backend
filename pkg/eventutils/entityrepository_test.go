package eventutils

import (
	"testing"
)

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
