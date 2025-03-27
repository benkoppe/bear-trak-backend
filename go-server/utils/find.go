package utils

func Find[T any](slice []T, predicate func(T) bool) *T {
	for i, item := range slice {
		if predicate(item) {
			return &slice[i]
		}
	}
	return nil
}

func Contains[T comparable](slice []T, target T) bool {
	return Find(slice, func(item T) bool {
		return item == target
	}) != nil
}
