package eventutils2

import (
	"github.com/gofrs/uuid"
)

type Entity struct {
	ID        uuid.UUID
	Events    []IEvent
	IsDeleted bool
}

func reconstructFromEvents(entity IReconstructible, events []Event) {
	for _, event := range events {
		entity.ApplyEvent(event)
	}
}
