package util

import (
	"fmt"
	"strconv"

	"golang.org/x/exp/constraints"
)

func StrToInt64(value string) int64 {
	if value, err := strconv.ParseInt(value, 10, 64); err == nil {
		return value
	}
	return 0
}

func IntToStr[T constraints.Integer](value T) string {
	return fmt.Sprintf("%d", value)
}
