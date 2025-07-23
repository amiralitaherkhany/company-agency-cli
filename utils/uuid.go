package utils

import (
	"log"

	"github.com/google/uuid"
)

func GenerateUUID() string {
	defer func() {
		if err := recover(); err != nil {
			log.Fatalln(err)
		}
	}()
	ID := uuid.New()
	return ID.String()
}
