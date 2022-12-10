package events

import (
	"github.com/Resta-Inc/resta/pkg/eventutils"
	"github.com/gofrs/uuid"
)

type CategoryCreated struct {
	eventutils.EventInfo
	Name         string
	ParentMenuID uuid.UUID
}

type CategoryNameChanged struct {
	eventutils.EventInfo
	NewName string
}

type SubCategoryAddedToCategory struct {
	eventutils.EventInfo
	SubCategoryID uuid.UUID
}
