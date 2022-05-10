package events

import (
	"github.com/Resta-Inc/resta/pkg/eventutils"
	"github.com/gofrs/uuid"
)

type SubCategoryCreated struct {
	eventutils.EntityEventInfo
	Name             string
	ParentCategoryID uuid.UUID
}

type SubCategoryNameChanged struct {
	eventutils.EntityEventInfo
	NewName string
}
