package eventutils2

import (
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/stretchr/testify/mock"
)

type MockEventStore struct {
	mock.Mock
}

func (m MockEventStore) SaveEventsToNewStream(streamName string, events []Event) (*esdb.WriteResult, error) {
	args := m.Called(streamName, events)
	writeResult, _ := args.Get(0).(*esdb.WriteResult)
	return writeResult, args.Error(1)
}
func (m MockEventStore) SaveEventsToExistingStream(streamName string, events []Event) (*esdb.WriteResult, error) {
	args := m.Called(streamName, events)
	writeResult, _ := args.Get(0).(*esdb.WriteResult)
	return writeResult, args.Error(1)
}

func (m MockEventStore) GetAllEventsByStreamName(streamName string) ([]Event, error) {
	args := m.Called(streamName)
	returnedEvents, _ := args.Get(0).([]Event)
	return returnedEvents, args.Error(1)
}

// type MockEntityRepository struct {
// 	mock.Mock
// }

// func (m MockEntityRepository) GetEntity(entity IEntity, id uuid.UUID) (IEntity, error) {
// 	args := m.Called(entity, id)
// 	returnedEntity, _ := args.Get(0).(IEntity)
// 	return returnedEntity, args.Error(1)
// }

// func (m MockEntityRepository) SaveEntity(entity IEntity) error {
// 	args := m.Called(entity)
// 	return args.Error(0)
// }
