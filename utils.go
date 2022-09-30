package n_error

// SliceContain if target in slice return true
func SliceContain[K comparable](target K, slice []K) bool {
	for _, a := range slice {
		if a == target {
			return true
		}
	}

	return false
}
