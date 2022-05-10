package events

import (
	"github.com/Resta-Inc/resta/pkg/eventutils"
	"github.com/gofrs/uuid"
)

type MenuItemCreated struct {
	eventutils.EntityEventInfo
	Name                string
	ParentSubCategoryID uuid.UUID
}
