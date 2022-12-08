package eventutils2

import (
	"context"
	"errors"
	"io"

	"github.com/EventStore/EventStore-Client-Go/esdb"
)

func NewEsdbClient(connectionString string) (*esdb.Client, error) {
	settings, err := esdb.ParseConnectionString(connectionString)
	if err != nil {
		return nil, err
	}
	db, err := esdb.NewClient(settings)
	return db, err
}

type IEventStore interface {
	SaveEventsToNewStream(streamName string, events []Event) (*esdb.WriteResult, error)
	SaveEventsToExistingStream(streamName string, events []Event) (*esdb.WriteResult, error)
	GetAllEventsByStreamName(streamName string) ([]Event, error)
}

type EventStore struct {
	db *esdb.Client
}

func NewEventStore(client *esdb.Client) (*EventStore, error) {
	eventStore := &EventStore{
		db: client,
	}
	return eventStore, nil
}

func (eventStore EventStore) SaveEventsToNewStream(streamName string, events []Event) (*esdb.WriteResult, error) {
	batch := prepareEventsBatch(events)
	options := esdb.AppendToStreamOptions{
		ExpectedRevision: esdb.NoStream{},
	}
	writeResult, err := eventStore.db.AppendToStream(context.Background(), streamName, options, batch...)
	return writeResult, err
}

func (eventStore EventStore) SaveEventsToExistingStream(streamName string, events []Event) (*esdb.WriteResult, error) {
	batch := prepareEventsBatch(events)
	options := esdb.AppendToStreamOptions{
		ExpectedRevision: esdb.StreamExists{},
	}
	writeResult, err := eventStore.db.AppendToStream(context.Background(), streamName, options, batch...)
	return writeResult, err
}

func (eventStore EventStore) GetAllEventsByStreamName(streamName string) ([]Event, error) {

	options := esdb.ReadStreamOptions{}
	stream, err := eventStore.db.ReadStream(context.Background(), streamName, options, 200)
	if errors.Is(err, esdb.ErrStreamNotFound) {
		return nil, ErrResourceNotFound
	}
	if err != nil {
		return nil, err
	}
	defer stream.Close()
	events := []Event{}
	for {
		eventData, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		event := Event{
			ID:   eventData.Event.EventID,
			Name: eventData.Event.EventType,
			Data: eventData.Event.Data,
		}
		events = append(events, event)
	}
	return events, nil
}

func DeserializeRecordedEvent(recordedEvent *esdb.RecordedEvent) Event {
	return Event{
		ID:   recordedEvent.EventID,
		Name: recordedEvent.EventType,
		Data: recordedEvent.Data,
	}
}

func prepareEventsBatch(events []Event) []esdb.EventData {
	batch := []esdb.EventData{}
	for _, event := range events {
		eventData := esdb.EventData{
			EventID:     event.ID,
			ContentType: esdb.JsonContentType,
			EventType:   event.Name,
			Data:        event.Data,
		}
		batch = append(batch, eventData)
	}
	return batch
}

// Errors

var (
	ErrResourceNotFound = errors.New("resource was not found")
)
