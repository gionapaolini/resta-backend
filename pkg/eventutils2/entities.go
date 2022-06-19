package eventutils2

import (
	"github.com/gofrs/uuid"
)

type Entity struct {
	ID        uuid.UUID
	Events    []IEvent
	IsDeleted bool
}
