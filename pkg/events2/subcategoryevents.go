package events2

import (
	"github.com/Resta-Inc/resta/pkg/eventutils2"
	"github.com/gofrs/uuid"
)

type SubCategoryCreated struct {
	eventutils2.EventInfo
	Name             string
	ParentCategoryID uuid.UUID
}

type SubCategoryNameChanged struct {
	eventutils2.EventInfo
	NewName string
}

type MenuItemAddedToSubCategory struct {
	eventutils2.EventInfo
	MenuItemID uuid.UUID
}
