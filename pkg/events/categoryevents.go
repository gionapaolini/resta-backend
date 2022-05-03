package events

import (
	"github.com/Resta-Inc/resta/pkg/eventutils"
	"github.com/gofrs/uuid"
)

type CategoryCreated struct {
	eventutils.EntityEventInfo
	Name         string
	ImageURL     string
	ParentMenuID uuid.UUID
}

type CategoryNameChanged struct {
	eventutils.EntityEventInfo
	NewName string
}
