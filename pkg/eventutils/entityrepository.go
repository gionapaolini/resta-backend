package eventutils

import (
	"errors"

	"github.com/Resta-Inc/resta/pkg/utils"
	"github.com/gofrs/uuid"
)

type IEntityRepository interface {
	GetEntity(entity IReconstructible, id uuid.UUID) (IReconstructible, error)
	SaveEntity(entity IReconstructible) error
}

type EntityRepository struct {
	EventStore IEventStore
}

func NewEntityRepository(eventStore IEventStore) IEntityRepository {
	return &EntityRepository{
		EventStore: eventStore,
	}
}

func (repo EntityRepository) GetEntity(entity IReconstructible, id uuid.UUID) (IReconstructible, error) {
	streamName := getStreamNameWithID(entity, id)
	returnedEvents, err := repo.EventStore.GetAllEventsByStreamName(streamName)
	if errors.Is(err, ErrResourceNotFound) {
		return nil, ErrEntityNotFound
	}
	ReconstructFromEvents(entity, returnedEvents)
	return entity, nil
}

func (repo EntityRepository) SaveEntity(entity IReconstructible) error {
	var err error
	if entity.IsNew() {
		_, err = repo.EventStore.SaveEventsToNewStream(getStreamName(entity), serializeEvents(entity.GetEvents()))
	} else {
		_, err = repo.EventStore.SaveEventsToExistingStream(getStreamName(entity), serializeEvents(entity.GetEvents()))
	}
	return err
}

func serializeEvents(events []IEvent) []Event {
	serializedEvents := []Event{}
	for _, event := range events {
		serializedEvents = append(serializedEvents, SerializedEvent(event))
	}
	return serializedEvents
}

func getStreamName(entity IReconstructible) string {
	return utils.GetType(entity) + "_" + entity.GetID().String()
}

func getStreamNameWithID(entity IReconstructible, id uuid.UUID) string {
	return utils.GetType(entity) + "_" + id.String()
}

// Errors

var (
	ErrEntityNotFound = errors.New("entity was not found")
)
