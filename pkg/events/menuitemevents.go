package events

import (
	"time"

	"github.com/Resta-Inc/resta/pkg/eventutils"
	"github.com/gofrs/uuid"
)

type MenuItemCreated struct {
	eventutils.EntityEventInfo
	Name                string
	ParentSubCategoryID uuid.UUID
}

type MenuItemNameChanged struct {
	eventutils.EntityEventInfo
	NewName string
}

type MenuItemEstimatedPreparationTimeChanged struct {
	eventutils.EntityEventInfo
	NewEstimate time.Duration
}
