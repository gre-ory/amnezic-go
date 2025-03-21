package util

import "strings"

// //////////////////////////////////////////////////
// join

func Join[T any](items []T, separator string) string {
	return ConvertAndJoin(items, ToStr[T], separator)
}

func JoinMap[K comparable, V any](items map[K]V, separator string) string {
	return ConvertAndJoinMap(items, PairToStr[K, V], separator)
}

func ConvertAndJoin[T any](items []T, convert func(item T) string, separator string) string {
	return strings.Join(Convert(items, convert), separator)
}

func ConvertAndJoinMap[K comparable, V any](items map[K]V, convert func(key K, value V) string, separator string) string {
	return strings.Join(ConvertMap(items, convert), separator)
}
