package eventutils

import (
	"context"
	"errors"
	"io"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/Resta-Inc/resta/pkg/utils"
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
	SaveEventsToNewStream(streamName string, events []IEvent) (*esdb.WriteResult, error)
	SaveEventsToExistingStream(streamName string, events []IEvent) (*esdb.WriteResult, error)
	GetAllEventsByStreamName(streamName string) ([]ReturnedEvent, error)
}

type EventStore struct {
	db *esdb.Client
}

func NewEventStore(connectionString string) (*EventStore, error) {
	client, err := NewEsdbClient(connectionString)
	if err != nil {
		return nil, err
	}
	eventStore := &EventStore{
		db: client,
	}
	return eventStore, nil
}

func (eventStore EventStore) SaveEventsToNewStream(streamName string, events []IEvent) (*esdb.WriteResult, error) {
	batch := prepareEventsBatch(events)
	options := esdb.AppendToStreamOptions{
		ExpectedRevision: esdb.NoStream{},
	}
	writeResult, err := eventStore.db.AppendToStream(context.Background(), streamName, options, batch...)
	return writeResult, err
}

func (eventStore EventStore) SaveEventsToExistingStream(streamName string, events []IEvent) (*esdb.WriteResult, error) {
	batch := prepareEventsBatch(events)
	options := esdb.AppendToStreamOptions{
		ExpectedRevision: esdb.StreamExists{},
	}
	writeResult, err := eventStore.db.AppendToStream(context.Background(), streamName, options, batch...)
	return writeResult, err
}

func (eventStore EventStore) GetAllEventsByStreamName(streamName string) ([]ReturnedEvent, error) {

	options := esdb.ReadStreamOptions{}
	stream, err := eventStore.db.ReadStream(context.Background(), streamName, options, 200)
	if errors.Is(err, esdb.ErrStreamNotFound) {
		return nil, ErrResourceNotFound
	}
	if err != nil {
		return nil, err
	}
	defer stream.Close()
	events := []ReturnedEvent{}
	for {
		eventData, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		event := ReturnedEvent{
			Data: eventData.Event.Data,
		}
		events = append(events, event)
	}
	return events, nil
}

func prepareEventsBatch(events []IEvent) []esdb.EventData {
	batch := []esdb.EventData{}
	for _, event := range events {
		finalJson := utils.SerializeObject(event)
		eventData := esdb.EventData{
			EventID:     event.GetEventID(),
			ContentType: esdb.JsonContentType,
			EventType:   utils.GetType(event),
			Data:        finalJson,
		}
		batch = append(batch, eventData)
	}
	return batch
}

// Errors

var (
	ErrResourceNotFound = errors.New("resource was not found")
)
