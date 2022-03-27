package eventutils

import (
	"encoding/json"
	"errors"

	"github.com/Resta-Inc/resta/pkg/utils"

	"github.com/gofrs/uuid"
)

type IEntity interface {
	GetID() uuid.UUID
	GetLatestEvents() []IEvent
	GetCommittedEvents() []IEvent
	GetAllEvents() []IEvent
	DeserializeEvent(jsonData []byte) IEvent

	isEmptyEntity() bool
	setLatestEvents(events []IEvent)
	setCommittedEvents(events []IEvent)
}

type Entity struct {
	ID              uuid.UUID
	CommittedEvents []IEvent
	LatestEvents    []IEvent
	IsDeleted       bool
}

func (entity Entity) GetID() uuid.UUID {
	return entity.ID
}

func (entity Entity) GetCommittedEvents() []IEvent {
	return entity.CommittedEvents
}

func (entity Entity) GetAllEvents() []IEvent {
	return append(entity.CommittedEvents, entity.LatestEvents...)
}

func (entity Entity) GetLatestEvents() []IEvent {
	return entity.LatestEvents
}

func (entity Entity) isEmptyEntity() bool {
	return len(entity.GetAllEvents()) == 0
}

func (entity *Entity) setLatestEvents(events []IEvent) {
	entity.LatestEvents = events
}

func (entity *Entity) setCommittedEvents(events []IEvent) {
	entity.CommittedEvents = events
}

type IEntityEvent interface {
	IEvent
	GetEntityID() uuid.UUID
	Apply(entity IEntity) IEntity
}

type EntityEventInfo struct {
	EventInfo
	EntityID uuid.UUID
}

func NewEntityEventInfo(entityID uuid.UUID) EntityEventInfo {
	entityEventInfo := EntityEventInfo{
		EventInfo: NewEventInfo(),
		EntityID:  entityID,
	}
	return entityEventInfo
}

func (entityEventInfo EntityEventInfo) GetEntityID() uuid.UUID {
	return entityEventInfo.EntityID
}

func ReconstructFromEvents(entity IEntity, events []IEvent) (IEntity, error) {
	if !entity.isEmptyEntity() {
		return nil, ErrEntityNotEmpty
	}
	for _, event := range events {
		entity = AddEvent(entity, event)
	}
	return entity, nil
}

func AddNewEvent(entity IEntity, event IEvent) IEntity {
	entity.setLatestEvents(append(entity.GetLatestEvents(), event))
	entity = ApplyEvent(entity, event.(IEntityEvent))
	return entity
}

func AddEvent(entity IEntity, event IEvent) IEntity {
	entity.setCommittedEvents(append(entity.GetCommittedEvents(), event))
	entity = ApplyEvent(entity, event.(IEntityEvent))
	return entity
}

func ApplyEvent(entity IEntity, event IEntityEvent) IEntity {
	return event.Apply(entity)
}

func GetRawDataFromSerializedEvent(jsonData []byte) (eventType string, rawData json.RawMessage) {
	var event utils.TypedJson
	json.Unmarshal(jsonData, &event)
	eventType, rawData = event.Type, event.Data
	return
}

func GetStreamName(entity IEntity) string {
	return utils.GetType(entity) + "_" + entity.GetID().String()
}

func GetStreamNameWithID(entity IEntity, id uuid.UUID) string {
	return utils.GetType(entity) + "_" + id.String()
}

func CommitEvents(entity IEntity) IEntity {
	committedEvents := entity.GetCommittedEvents()
	entity.setCommittedEvents(append(committedEvents, entity.GetLatestEvents()...))
	entity.setLatestEvents(nil)
	return entity
}

// Errors
var (
	ErrEntityNotEmpty = errors.New("entity is not empty")
)
