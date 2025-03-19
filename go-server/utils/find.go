package utils

func Find[T any](slice []T, predicate func(T) bool) *T {
	for i, item := range slice {
		if predicate(item) {
			return &slice[i]
		}
	}
	return nil
}
