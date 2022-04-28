package events

import (
	"github.com/Resta-Inc/resta/pkg/eventutils"
)

type CategoryCreated struct {
	eventutils.EntityEventInfo
	Name     string
	ImageURL string
}
