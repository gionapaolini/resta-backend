package events

import (
	"time"

	"github.com/Resta-Inc/resta/pkg/eventutils"
	"github.com/gofrs/uuid"
)

type MenuItemCreated struct {
	eventutils.EventInfo
	Name                string
	ParentSubCategoryID uuid.UUID
}

type MenuItemNameChanged struct {
	eventutils.EventInfo
	NewName string
}

type MenuItemEstimatedPreparationTimeChanged struct {
	eventutils.EventInfo
	NewEstimate time.Duration
}
