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

type IEventInfo interface {
	GetEventID() uuid.UUID
}

type EventInfo struct {
	EntityID, EventID uuid.UUID
	CreatedAt         time.Time
}

func (ei EventInfo) GetEventID() uuid.UUID {
	return ei.EventID
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
}

type IEvent interface {
	GetEventID() uuid.UUID
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
