package util

// //////////////////////////////////////////////////
// find if

func FindIf[T any](items []T, predicate func(item T) bool) (T, bool) {
	for _, item := range items {
		if predicate(item) {
			return item, true
		}
	}
	var empty T
	return empty, false
}
