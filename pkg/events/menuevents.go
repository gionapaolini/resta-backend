package events

import (
	"github.com/Resta-Inc/resta/pkg/eventutils"
	"github.com/gofrs/uuid"
)

type MenuCreated struct {
	eventutils.EventInfo
	Name string
}

type MenuEnabled struct {
	eventutils.EventInfo
}

type MenuDisabled struct {
	eventutils.EventInfo
}

type MenuNameChanged struct {
	eventutils.EventInfo
	NewName string
}

type CategoryAddedToMenu struct {
	eventutils.EventInfo
	CategoryID uuid.UUID
}
