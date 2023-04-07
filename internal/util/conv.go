package util

import "strconv"

func StrToInt64(value string) int64 {
	if value, err := strconv.ParseInt(value, 10, 64); err == nil {
		return value
	}
	return 0
}
