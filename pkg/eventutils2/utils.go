package eventutils2

import (
	"encoding/json"
	"time"

	"github.com/Resta-Inc/resta/pkg/utils"
	"github.com/gofrs/uuid"
)

type Event struct {
	ID   uuid.UUID
	Name string
	Data []byte
}

type EventInfo struct {
	EntityID, EventID uuid.UUID
	CreatedAt         time.Time
}

func (ei EventInfo) GetEventID() uuid.UUID {
	return ei.EventID
}

func (ei EventInfo) GetEntityID() uuid.UUID {
	return ei.EntityID
}

func (ei EventInfo) GetTimeStamp() time.Time {
	return ei.CreatedAt
}

func NewEventInfo(entityID uuid.UUID) EventInfo {
	eventID := utils.GenerateNewUUID()
	eventInfo := EventInfo{
		EventID:   eventID,
		EntityID:  entityID,
		CreatedAt: utils.Time.Now(),
	}
	return eventInfo
}

type Entity struct {
	ID        uuid.UUID
	Events    []IEvent
	IsDeleted bool
	New       bool
}

func (e Entity) GetEvents() []IEvent {
	return e.Events
}

func (e Entity) GetID() uuid.UUID {
	return e.ID
}

func (e *Entity) SetNew() {
	e.New = true
}

func (e Entity) IsNew() bool {
	return e.New
}

type IReconstructible interface {
	GetID() uuid.UUID
	GetEvents() []IEvent
	DeserializeEvent(event Event) IEvent
	ApplyEvent(event IEvent)
	AppendEvent(event IEvent)
	IsNew() bool
	SetNew()
}

type IEvent interface {
	GetEventID() uuid.UUID
	GetEntityID() uuid.UUID
	GetTimeStamp() time.Time
}

func ReconstructFromEvents(entity IReconstructible, events []Event) {
	for _, event := range events {
		deserializedEvent := entity.DeserializeEvent(event)
		entity.ApplyEvent(deserializedEvent)
	}
}

func SerializedEvent(eventObj IEvent) Event {
	bytes, _ := json.Marshal(eventObj)
	event := Event{
		ID:   eventObj.GetEventID(),
		Name: utils.GetType(eventObj),
		Data: bytes,
	}
	return event
}

func AddEvent(event IEvent, obj IReconstructible) {
	obj.ApplyEvent(event)
	obj.AppendEvent(event)
}
