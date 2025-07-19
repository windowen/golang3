package utils

// SliceHas 是否在切片
func SliceHas[T comparable](slice []T, val T) bool {
	for _, ele := range slice {
		if ele == val {
			return true
		}
	}

	return false
}
