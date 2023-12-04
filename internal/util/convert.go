package util

// //////////////////////////////////////////////////
// convert

func Convert[T any, U any](items []T, convert func(item T) U) []U {
	converted := make([]U, 0, len(items))
	for _, item := range items {
		converted = append(converted, convert(item))
	}
	return converted
}

func ConvertToMap[T any, K comparable, V any](items []T, convert func(item T) (K, V)) map[K]V {
	converted := make(map[K]V, len(items))
	for _, item := range items {
		key, value := convert(item)
		converted[key] = value
	}
	return converted
}

func ConvertMap[K comparable, V any, U any](items map[K]V, convert func(key K, value V) U) []U {
	converted := make([]U, 0, len(items))
	for key, value := range items {
		converted = append(converted, convert(key, value))
	}
	return converted
}
