// //go:build integration

package eventutils2

import (
	"testing"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/Resta-Inc/resta/pkg/utils"
	"github.com/stretchr/testify/require"
)

const eventStoreConnectionString = "esdb://127.0.0.1:2113?tls=false&keepAliveTimeout=10000&keepAliveInterval=10000"

type TestEvent struct {
	EventInfo
	Data string
}

func TestSaveEventsAsNewStream(t *testing.T) {
	// Arrange
	settings, _ := esdb.ParseConnectionString(eventStoreConnectionString)
	db, _ := esdb.NewClient(settings)
	eventStore, err := NewEventStore(db)

	require.NoError(t, err)

	events := getRandomEvents(5)

	streamName := "TestStream-" + utils.GenerateNewUUID().String()

	// Act
	_, err = eventStore.SaveEventsToNewStream(streamName, events)
	require.NoError(t, err)
	returnedEvents, err := eventStore.GetAllEventsByStreamName(streamName)
	require.NoError(t, err)

	// Assert
	require.Equal(t, returnedEvents, events)
}

func TestSaveEventsToExistentStream(t *testing.T) {
	// Arrange
	settings, _ := esdb.ParseConnectionString(eventStoreConnectionString)
	db, _ := esdb.NewClient(settings)
	eventStore, err := NewEventStore(db)

	events := getRandomEvents(5)
	newEvents := getRandomEvents(5)

	streamName := "TestStream-" + utils.GenerateNewUUID().String()
	_, _ = eventStore.SaveEventsToNewStream(streamName, events)

	// Act
	_, err = eventStore.SaveEventsToExistingStream(streamName, newEvents)
	require.NoError(t, err)

	// Assert
	returnedEvents, err := eventStore.GetAllEventsByStreamName(streamName)
	require.NoError(t, err)
	require.Equal(t, returnedEvents, append(events, newEvents...))
}

func getRandomEvents(n int) []Event {
	events := []Event{}
	for i := 0; i < n; i++ {
		event := TestEvent{
			EventInfo: NewEventInfo(utils.GenerateNewUUID()),
		}
		event.Data = "Test - " + event.GetEventID().String()
		serializedEvent := SerializedEvent(event)
		events = append(events, serializedEvent)
	}
	return events
}
