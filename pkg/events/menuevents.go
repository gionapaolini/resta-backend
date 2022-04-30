package events

import (
	"github.com/Resta-Inc/resta/pkg/eventutils"
	"github.com/gofrs/uuid"
)

type MenuCreated struct {
	eventutils.EntityEventInfo
	Name string
}

type MenuEnabled struct {
	eventutils.EntityEventInfo
}

type MenuDisabled struct {
	eventutils.EntityEventInfo
}

type MenuNameChanged struct {
	eventutils.EntityEventInfo
	NewName string
}

type CategoryAddedToMenu struct {
	eventutils.EntityEventInfo
	CategoryID uuid.UUID
}
