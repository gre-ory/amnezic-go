package util

// //////////////////////////////////////////////////
// compare

func EqualMap[K comparable, V comparable](left, right map[K]V) bool {
	if left == nil && right == nil {
		return true
	}
	if left == nil || right == nil {
		return false
	}
	if len(left) != len(right) {
		return false
	}
	for key, leftValue := range left {
		rightValue, found := right[key]
		if !found {
			return false
		}
		if leftValue != rightValue {
			return false
		}
	}
	for key := range right {
		_, found := left[key]
		if !found {
			return false
		}
	}
	return true
}
