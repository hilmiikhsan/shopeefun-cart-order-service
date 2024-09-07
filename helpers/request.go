package helpers

import (
	"fmt"
	"time"
)

func GenerateRefCode() string {
	currentUnixTime := time.Now().Unix()

	refCode := fmt.Sprintf("REF%d", currentUnixTime)

	return refCode
}
