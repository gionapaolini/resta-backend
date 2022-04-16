// //go:build integration

package eventutils

import (
	"encoding/json"
	"testing"

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
	eventStore, err := NewEventStore(eventStoreConnectionString)

	require.NoError(t, err)

	events := getRandomEvents(5)

	streamName := "TestStream-" + utils.GenerateNewUUID().String()

	// Act
	_, err = eventStore.SaveEventsToNewStream(streamName, events)
	require.NoError(t, err)
	returnedEvents, err := eventStore.GetAllEventsByStreamName(streamName)
	require.NoError(t, err)

	// Assert
	storedEvents := deserializeEvents(returnedEvents)
	require.Equal(t, storedEvents, events)
}

func TestSaveEventsToExistentStream(t *testing.T) {
	// Arrange
	eventStore, _ := NewEventStore(eventStoreConnectionString)
	events := getRandomEvents(5)
	newEvents := getRandomEvents(5)

	streamName := "TestStream-" + utils.GenerateNewUUID().String()
	_, _ = eventStore.SaveEventsToNewStream(streamName, events)

	// Act
	_, err := eventStore.SaveEventsToExistingStream(streamName, newEvents)
	require.NoError(t, err)

	// Assert
	returnedEvents, err := eventStore.GetAllEventsByStreamName(streamName)
	require.NoError(t, err)
	storedEvents := deserializeEvents(returnedEvents)
	require.Equal(t, storedEvents, append(events, newEvents...))
}

func getRandomEvents(n int) []IEvent {
	events := []IEvent{}
	for i := 0; i < n; i++ {
		event := TestEvent{
			EventInfo: NewEventInfo(),
		}
		event.Data = "Test - " + event.GetEventID().String()
		events = append(events, event)
	}
	return events
}

func deserializeEvents(returnedEvents []ReturnedEvent) []IEvent {
	storedEvents := []IEvent{}
	for _, returnedEvent := range returnedEvents {
		var event utils.TypedJson
		json.Unmarshal(returnedEvent.Data, &event)
		var deserializedEvent TestEvent
		json.Unmarshal(event.Data, &deserializedEvent)
		storedEvents = append(storedEvents, deserializedEvent)
	}
	return storedEvents
}
