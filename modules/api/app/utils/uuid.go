package utils

import (
	"strings"

	"github.com/satori/go.uuid"
)

func GenerateUUID() string {
	sig := uuid.NewV1().String()
	sig = strings.Replace(sig, "-", "", -1)
	return sig
}
