package utils

func KeysIntersection[M1 ~map[K]V1, M2 ~map[K]V2, K comparable, V1 any, V2 any](left M1, right M2) []K {
	var keys []K

	for k := range left {
		if _, found := right[k]; found {
			keys = append(keys, k)
		}
	}

	return keys
}
