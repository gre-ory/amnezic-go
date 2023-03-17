package util

import "strings"

// //////////////////////////////////////////////////
// join

func Join[T any](items []T, convert func(item T) string) string {
	return strings.Join(Convert(items, convert), ",")
}
