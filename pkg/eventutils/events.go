package eventutils

import (
	"time"

	"github.com/Resta-Inc/resta/pkg/utils"
	"github.com/gofrs/uuid"
)

type IEvent interface {
	GetEventID() uuid.UUID
	GetDateTime() time.Time
}

type EventInfo struct {
	DateTime time.Time
	EventID  uuid.UUID
}

func NewEventInfo() EventInfo {
	eventID := utils.GenerateNewUUID()
	eventInfo := EventInfo{
		EventID:  eventID,
		DateTime: utils.Time.Now(),
	}
	return eventInfo
}

func (eventInfo EventInfo) GetDateTime() time.Time {
	return eventInfo.DateTime
}

func (eventInfo EventInfo) GetEventID() uuid.UUID {
	return eventInfo.EventID
}

type ReturnedEvent struct {
	Data []byte
}
