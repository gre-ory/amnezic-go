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
