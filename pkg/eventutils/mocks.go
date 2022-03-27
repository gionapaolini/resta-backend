package eventutils

import (
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/stretchr/testify/mock"
)

type MockEventStore struct {
	mock.Mock
}

func (m MockEventStore) SaveEventsToNewStream(streamName string, events []IEvent) (*esdb.WriteResult, error) {
	args := m.Called(streamName, events)
	writeResult, _ := args.Get(0).(*esdb.WriteResult)
	return writeResult, args.Error(1)
}
func (m MockEventStore) SaveEventsToExistingStream(streamName string, events []IEvent) (*esdb.WriteResult, error) {
	args := m.Called(streamName, events)
	writeResult, _ := args.Get(0).(*esdb.WriteResult)
	return writeResult, args.Error(1)
}

func (m MockEventStore) GetAllEventsByStreamName(streamName string) ([]ReturnedEvent, error) {
	args := m.Called(streamName)
	returnedEvents, _ := args.Get(0).([]ReturnedEvent)
	return returnedEvents, args.Error(1)
}
