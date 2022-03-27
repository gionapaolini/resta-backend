package eventutils

import "github.com/gofrs/uuid"

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
