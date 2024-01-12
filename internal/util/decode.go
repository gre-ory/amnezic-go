package util

// //////////////////////////////////////////////////
// decode

func Decode[T any, U any](items []T, decode func(item T) (U, error)) []U {
	decoded := make([]U, 0, len(items))
	for _, item := range items {
		value, err := decode(item)
		if err != nil {
			decoded = append(decoded, value)
		}
	}
	return decoded
}
