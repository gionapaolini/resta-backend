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

type IReconstructible interface {
	ApplyEvent(event Event)
	AppendEvent(event IEvent)
}

type IEvent interface {
	GetEventID() uuid.UUID
	GetEntityID() uuid.UUID
	GetTimeStamp() time.Time
}

func ReconstructFromEvents(entity IReconstructible, events []Event) {
	for _, event := range events {
		entity.ApplyEvent(event)
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
	serializedEvent := SerializedEvent(event)
	obj.ApplyEvent(serializedEvent)
	obj.AppendEvent(event)
}
