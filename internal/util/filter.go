package util

// //////////////////////////////////////////////////
// filter

func Filter[T any](items []T, predicate func(item T) bool) []T {
	filtered := make([]T, 0, len(items))
	for _, item := range items {
		if predicate(item) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}
