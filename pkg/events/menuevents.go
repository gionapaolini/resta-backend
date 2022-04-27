package events

import (
	"github.com/Resta-Inc/resta/pkg/eventutils"
)

type MenuCreated struct {
	eventutils.EntityEventInfo
	Name string
}
