package events2

import (
	"time"

	"github.com/Resta-Inc/resta/pkg/eventutils2"
	"github.com/gofrs/uuid"
)

type MenuItemCreated struct {
	eventutils2.EventInfo
	Name                string
	ParentSubCategoryID uuid.UUID
}

type MenuItemNameChanged struct {
	eventutils2.EventInfo
	NewName string
}

type MenuItemEstimatedPreparationTimeChanged struct {
	eventutils2.EventInfo
	NewEstimate time.Duration
}
