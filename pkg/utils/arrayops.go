package utils

import "github.com/gofrs/uuid"

func FindID(array []uuid.UUID, item uuid.UUID) int {
	for i, v := range array {
		if v == item {
			return i
		}
	}
	return -1
}

func RemoveID(slice []uuid.UUID, s int) []uuid.UUID {
	return append(slice[:s], slice[s+1:]...)
}
