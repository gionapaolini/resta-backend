package eventutils

import (
	"errors"

	"github.com/gofrs/uuid"
)

type IEntityRepository interface {
	GetEntity(entity IEntity, id uuid.UUID) (IEntity, error)
	SaveEntity(entity IEntity) error
}

type EntityRepository struct {
	EventStore IEventStore
}

func NewEntityRepository(eventStore IEventStore) IEntityRepository {
	return &EntityRepository{
		EventStore: eventStore,
	}
}

func (repo EntityRepository) GetEntity(entity IEntity, id uuid.UUID) (IEntity, error) {
	streamName := GetStreamNameWithID(entity, id)
	returnedEvents, err := repo.EventStore.GetAllEventsByStreamName(streamName)
	if errors.Is(err, ErrResourceNotFound) {
		return nil, ErrEntityNotFound
	}
	var events []IEvent
	for _, returnedEvent := range returnedEvents {
		deserializedEvent := entity.DeserializeEvent(returnedEvent.Data)
		events = append(events, deserializedEvent)
	}

	entity, err = ReconstructFromEvents(entity, events)
	if err != nil {
		panic(err)
	}

	return entity, nil
}

func (repo EntityRepository) SaveEntity(entity IEntity) error {
	if len(entity.GetCommittedEvents()) == 0 {
		repo.EventStore.SaveEventsToNewStream(GetStreamName(entity), entity.GetLatestEvents())
	} else {
		repo.EventStore.SaveEventsToExistingStream(GetStreamName(entity), entity.GetLatestEvents())
	}
	return nil
}

// Errors

var (
	ErrEntityNotFound = errors.New("entity was not found")
)
