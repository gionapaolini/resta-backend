package events2

import (
	"github.com/Resta-Inc/resta/pkg/eventutils2"
	"github.com/gofrs/uuid"
)

type MenuCreated struct {
	eventutils2.EventInfo
	Name string
}

type MenuEnabled struct {
	eventutils2.EventInfo
}

type MenuDisabled struct {
	eventutils2.EventInfo
}

type MenuNameChanged struct {
	eventutils2.EventInfo
	NewName string
}

type CategoryAddedToMenu struct {
	eventutils2.EventInfo
	CategoryID uuid.UUID
}
