package utils

func List2Map[K comparable, V any](list []V, keyFunc func(V) K) map[K]V {
	result := make(map[K]V, len(list))
	for _, item := range list {
		key := keyFunc(item)
		result[key] = item
	}
	return result
}

func Map[T any](list []T, mapFunc func(T) T) []T {
	result := make([]T, len(list))
	for i, item := range list {
		result[i] = mapFunc(item)
	}
	return result
}

func Filter[T any](list []T, filterFunc func(T, int) bool) []T {
	var result []T
	for index, item := range list {
		if filterFunc(item, index) {
			result = append(result, item)
		}
	}
	return result
}

func FindIndex[T any](list []T, findFunc func(T) bool) int {
	for i, item := range list {
		if findFunc(item) {
			return i
		}
	}
	return -1
}
