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

func ToStr[T any](value T) string {
	if stringer, ok := any(value).(fmt.Stringer); ok {
		return stringer.String()
	}
	return fmt.Sprintf("%v", value)
}

func PairToStr[K any, V any](key K, value V) string {
	return fmt.Sprintf("%s:%s", ToStr[K](key), ToStr[V](value))
}
