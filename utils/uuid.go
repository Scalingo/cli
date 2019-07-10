package utils

import (
	"github.com/google/uuid"
)

func IsUUID(str string) bool {
	_, err := uuid.Parse(str)

	return err == nil
}
