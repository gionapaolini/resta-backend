package events2

import (
	"github.com/Resta-Inc/resta/pkg/eventutils2"
	"github.com/gofrs/uuid"
)

type CategoryCreated struct {
	eventutils2.EventInfo
	Name         string
	ParentMenuID uuid.UUID
}

type CategoryNameChanged struct {
	eventutils2.EventInfo
	NewName string
}

type SubCategoryAddedToCategory struct {
	eventutils2.EventInfo
	SubCategoryID uuid.UUID
}
