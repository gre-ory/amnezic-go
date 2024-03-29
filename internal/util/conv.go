package util

import (
	"fmt"
	"strconv"
	"time"

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

func ToAny[T any](value T) any {
	return value
}

func TsToStr(ts int64) string {
	return time.Unix(ts, 0).Format(time.DateTime)
}
