package events

import (
	"github.com/Resta-Inc/resta/pkg/eventutils"
	"github.com/gofrs/uuid"
)

type CategoryCreated struct {
	eventutils.EntityEventInfo
	Name         string
	ParentMenuID uuid.UUID
}

type CategoryNameChanged struct {
	eventutils.EntityEventInfo
	NewName string
}

type SubCategoryAddedToCategory struct {
	eventutils.EntityEventInfo
	SubCategoryID uuid.UUID
}

type SubCategoryCreated struct {
	eventutils.EntityEventInfo
	Name             string
	ParentCategoryID uuid.UUID
}

type SubCategoryNameChanged struct {
	eventutils.EntityEventInfo
	NewName string
}
