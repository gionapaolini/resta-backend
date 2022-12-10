package events

import (
	"github.com/Resta-Inc/resta/pkg/eventutils"
	"github.com/gofrs/uuid"
)

type SubCategoryCreated struct {
	eventutils.EventInfo
	Name             string
	ParentCategoryID uuid.UUID
}

type SubCategoryNameChanged struct {
	eventutils.EventInfo
	NewName string
}

type MenuItemAddedToSubCategory struct {
	eventutils.EventInfo
	MenuItemID uuid.UUID
}
