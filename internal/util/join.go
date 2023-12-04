package util

import "strings"

// //////////////////////////////////////////////////
// join

func Join[T any](items []T, convert func(item T) string) string {
	return strings.Join(Convert(items, convert), ",")
}

func JoinMap[K comparable, V any](items map[K]V, convert func(key K, value V) string) string {
	return strings.Join(ConvertMap(items, convert), ",")
}
