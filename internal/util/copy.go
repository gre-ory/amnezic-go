package util

// //////////////////////////////////////////////////
// copy

func CopyMap[K comparable, V any](values map[K]V) map[K]V {
	if values == nil {
		return nil
	}
	copy := make(map[K]V, len(values))
	for key, value := range values {
		copy[key] = value
	}
	return copy
}
