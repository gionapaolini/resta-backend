package utils

import "github.com/gofrs/uuid"

func GenerateNewUUID() uuid.UUID {
	id, _ := uuid.NewV4()
	return id
}
